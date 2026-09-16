/* Deterministic, source-backed motion compositions. Loaded only by render.mjs. */
window.createFilmRenderer = async function (images, evidence) {
  const loaded = {}
  for (const [name, src] of Object.entries(images)) {
    const image = new Image()
    image.src = src
    await image.decode()
    loaded[name] = image
  }
  const canvas = document.createElement('canvas')
  document.body.append(canvas)
  const c = canvas.getContext('2d')
  let W,
    H,
    P,
    accent,
    ink = '#172a38'
  const ease = (t) => 1 - Math.pow(1 - Math.min(1, Math.max(0, t)), 3)
  const rect = (x, y, w, h, color, r = 14, stroke) => {
    c.beginPath()
    c.roundRect(x, y, w, h, r)
    c.fillStyle = color
    c.fill()
    if (stroke) {
      c.strokeStyle = stroke
      c.lineWidth = 2
      c.stroke()
    }
  }
  const text = (
    value,
    x,
    y,
    size = 30,
    color = ink,
    weight = 400,
    mono = false,
  ) => {
    c.fillStyle = color
    c.font = `${weight} ${size}px ${mono ? 'monospace' : 'Arial, sans-serif'}`
    c.textBaseline = 'top'
    c.fillText(value, x, y)
  }
  const lines = (
    value,
    x,
    y,
    width,
    size = 30,
    color = ink,
    weight = 400,
    gap = 1.35,
  ) => {
    c.font = `${weight} ${size}px Arial, sans-serif`
    const words = value.split(' ')
    let line = '',
      row = 0
    for (const word of words) {
      if (c.measureText(`${line} ${word}`.trim()).width > width && line) {
        text(line, x, y + row * size * gap, size, color, weight)
        line = word
        row++
      } else line = (line + ' ' + word).trim()
    }
    if (line) text(line, x, y + row * size * gap, size, color, weight)
    return (row + 1) * size * gap
  }
  const node = (x, y, w, h, title, detail, active = false) => {
    rect(
      x,
      y,
      w,
      h,
      active ? '#203b4a' : '#f9fbfc',
      14,
      active ? '#203b4a' : '#c0ccd3',
    )
    lines(
      title,
      x + 22,
      y + 18,
      w - 44,
      P ? 34 : 29,
      active ? '#f5f9fc' : ink,
      600,
      1.2,
    )
    if (detail)
      lines(
        detail,
        x + 22,
        y + h - 44,
        w - 44,
        P ? 27 : 21,
        active ? '#bfd9d6' : '#526b78',
      )
  }
  const edge = (x1, y1, x2, y2, t) => {
    c.beginPath()
    c.moveTo(x1, y1)
    c.lineTo(x2, y2)
    c.strokeStyle = '#aebfc8'
    c.lineWidth = 3
    c.stroke()
    const v = (t % 2.7) / 2.7
    c.beginPath()
    c.arc(x1 + (x2 - x1) * v, y1 + (y2 - y1) * v, 7, 0, Math.PI * 2)
    c.fillStyle = '#377f77'
    c.fill()
  }
  const flow = (items, t) => {
    if (P) {
      const x = 100,
        w = W - 200,
        h = 110,
        start = 292
      items.forEach((item, i) => {
        if (i)
          edge(
            W / 2,
            start + (i - 1) * 142 + h,
            W / 2,
            start + i * 142,
            t + i * 0.3,
          )
        node(
          x,
          start + i * 142,
          w,
          h,
          item[0],
          item[1],
          i === Math.min(items.length - 1, Math.floor(t / 2)),
        )
      })
    } else {
      const n = items.length,
        w = (W - 180 - (n - 1) * 44) / n,
        y = 435
      items.forEach((item, i) => {
        const x = 90 + i * (w + 44)
        if (i) edge(x - 44, y + 74, x, y + 74, t + i * 0.3)
        node(
          x,
          y,
          w,
          150,
          item[0],
          item[1],
          i === Math.min(n - 1, Math.floor(t / 2)),
        )
      })
    }
  }
  const headings = {
    raizhost: [
      'A website with a purpose.',
      'An edit with clear boundaries.',
      'A change with a path to production.',
    ],
    cloudcostmcp: [
      'Start with the infrastructure.',
      'Make the estimate inspectable.',
      'Keep the assumptions visible.',
    ],
    tickethacker: [
      'Every channel. One queue.',
      'Bring important work into focus.',
      'Keep the queue moving.',
    ],
    'llm-lint': [
      'Find the boundary violation.',
      'Explain what needs attention.',
      'See the fix before applying it.',
    ],
    'sre-reference-app': [
      'Reliability starts with the design.',
      'One task stops. Traffic continues.',
      'Replace. Recover. Measure.',
    ],
    'sre-landing-zone': [
      'Five accounts. Clear boundaries.',
      'Prepare the second region.',
      'Stop paying for idle development.',
    ],
  }
  const captions = {
    raizhost: [
      'A real public site, with a clear route from interest to inquiry.',
      'Owners edit bounded content. The site design stays in source.',
      'Content commit → site workflow → S3 → CloudFront.',
    ],
    cloudcostmcp: [
      'Parse the sample Terraform file into a resource inventory.',
      'A monthly estimate, broken down by resource and service.',
      'A live compute price and an estimated egress allowance are distinct.',
    ],
    tickethacker: [
      'The real agent inbox, populated with fictional sample tickets.',
      'Filter by priority to isolate the next ticket.',
      'Select a sample ticket and update its status from the inbox.',
    ],
    'llm-lint': [
      'A real scan flags a local configuration file in a sample repository.',
      'The finding includes a rule, severity, location, and remediation.',
      'Preview reports the proposed changes without applying them.',
    ],
    'sre-reference-app': [
      'The load balancer routes requests to two Fargate tasks.',
      'In the recorded exercise, the surviving task absorbed traffic.',
      'Historical observation: 78-second task recovery. Timeline compressed.',
    ],
    'sre-landing-zone': [
      'Separate workloads from the accounts responsible for audit and logs.',
      'Pilot light keeps standby compute at zero until it is needed.',
      'The scheduled role can stop only resources tagged Environment=dev.',
    ],
  }
  function terminal(id, scene, t) {
    let code,
      cmd,
      label = 'OUTPUT EXCERPT'
    if (id === 'cloudcostmcp') {
      const result = JSON.parse(evidence.cost.estimate.stdout)
      const compute = result.by_resource.find(
        (r) => r.resource_id === 'aws_instance.web',
      )
      if (scene === 0) {
        cmd = 'cloudcost analyze main.tf --json'
        code = [
          '"id": "aws_instance.web",',
          '"provider": "aws",',
          '"region": "us-east-1",',
          '"instance_type": "t3.small"',
          '',
          '"total_count": 1',
          '"parse_warnings": []',
        ]
      }
      if (scene === 1) {
        cmd = 'cloudcost estimate main.tf --json'
        code = [
          `"total_monthly": ${result.total_monthly},`,
          '"currency": "USD",',
          '"by_service": {',
          `  "compute": ${result.by_service.compute},`,
          `  "data_transfer": ${result.by_service.data_transfer}`,
          '}',
          '',
          'Sample: us-east-1 · 730h/month',
        ]
      }
      if (scene === 2) {
        cmd = 'Inspect pricing and warnings'
        label = 'REPORT EXCERPT'
        code = [
          `Compute pricing: ${compute.pricing_source}`,
          `Compute confidence: ${compute.confidence}`,
          '',
          'Estimated egress: 100 GB/month',
          `Egress allowance: $${result.estimated_egress_monthly.toFixed(2)}/month`,
          'Egress confidence: low',
          '',
          'Captured: 2026-09-16',
        ]
      }
    } else {
      if (scene === 0) {
        cmd = 'llm-lint scan'
        code = [
          '✕ LLM006  error',
          '.cursorrules / .cursor/ tracked',
          '',
          '  └─ .cursorrules',
          '',
          '1 findings',
          '1 errors, 0 warnings, 0 info',
        ]
      }
      if (scene === 1) {
        cmd = 'llm-lint rules show LLM006'
        code = [
          'ID:        LLM006',
          'Severity:  error',
          'Category:  cursor',
          'Kind:      path',
          '',
          'Remediation:',
          'Add local tool paths to .gitignore.',
          'Untrack existing entries.',
        ]
      }
      if (scene === 2) {
        cmd = 'llm-lint scan --fix-preview'
        code = [
          'would fix:',
          '  1 files changed',
          '  3 .gitignore entries added',
          '  1 index entries untracked',
          '',
          '0 commit messages cleaned',
          '',
          'Preview only. No changes applied.',
        ]
      }
    }
    const x = P ? 48 : 80,
      y = P ? 318 : 290,
      w = W - x * 2,
      h = P ? 570 : 495
    rect(x, y, w, h, '#10232d', 16)
    text(label, x + 32, y + 25, P ? 19 : 19, '#91a9b8', 500, true)
    c.strokeStyle = '#314653'
    c.lineWidth = 1
    c.beginPath()
    c.moveTo(x + 28, y + 65)
    c.lineTo(x + w - 28, y + 65)
    c.stroke()
    const size = P ? 28 : 28
    text('$', x + 30, y + 91, size, '#b1d9c1', 500, true)
    text(cmd, x + 62, y + 91, P ? 31 : 28, '#e5eff4', 500, true)
    const count = Math.min(code.length, Math.floor(t * 5) + 1)
    for (let i = 0; i < count; i++)
      text(
        code[i],
        x + 32,
        y + 157 + i * (P ? 44 : 36),
        P ? 36 : 27,
        i === 0 ? '#b8dfc6' : '#c5d3dd',
        400,
        true,
      )
  }
  function screen(id, scene, t) {
    const key =
      id === 'raizhost'
        ? P
          ? 'raizhost-portrait'
          : 'raizhost-wide'
        : ['ticket-inbox', 'ticket-context', 'ticket-resolved'][scene]
    const img = loaded[key]
    const x = P ? 48 : 72,
      y = P ? 280 : 250,
      w = W - 2 * x,
      h = P ? 595 : 550
    rect(x - 1, y - 1, w + 2, h + 2, '#bdcbd2', 17)
    c.save()
    c.beginPath()
    c.roundRect(x, y, w, h, 16)
    c.clip()
    if (id === 'raizhost' && P) {
      c.drawImage(img, 0, 180, img.width, (img.width * h) / w, x, y, w, h)
    } else if (id === 'tickethacker' && P) {
      // Readable crops of the real UI. The filter and result are separate details.
      if (scene === 1) {
        c.fillStyle = '#f8fafc'
        c.fillRect(x, y, w, h)
        c.drawImage(
          img,
          600,
          126,
          300,
          55,
          x + 24,
          y + 24,
          w - 48,
          ((w - 48) * 55) / 300,
        )
        c.drawImage(
          img,
          330,
          252,
          300,
          110,
          x + 24,
          y + 218,
          w - 48,
          ((w - 48) * 110) / 300,
        )
      } else c.drawImage(img, 328, 252, 420, (420 * h) / w, x, y, w, h)
    } else if (id === 'tickethacker')
      c.drawImage(img, 255, 60, 1025, (1025 * h) / w, x, y, w, h)
    else c.drawImage(img, 0, 0, img.width, (img.width * h) / w, x, y, w, h)
    c.restore()
    if (id === 'tickethacker') {
      rect(x + 18, y + h - 58, P ? 325 : 345, 42, '#142b36', 7)
      text(
        'SAMPLE WORKSPACE',
        x + 34,
        y + h - 47,
        P ? 21 : 22,
        '#fff',
        500,
        true,
      )
    }
  }
  function recovery(scene, t) {
    const left = P ? 75 : 90,
      top = P ? 307 : 360,
      w = P ? 750 : 1420
    const lbx = P ? 210 : 120,
      lby = P ? 300 : 422,
      lbw = P ? 480 : 320
    node(lbx, lby, lbw, 125, 'Load balancer', 'HTTP requests', true)
    const tx = P ? 100 : 660,
      ty = P ? 500 : 340,
      tw = P ? 325 : 335,
      th = 135
    if (P) {
      edge(W / 2, lby + 125, tx + tw / 2, ty, t)
      edge(W / 2, lby + 125, W - tx - tw / 2, ty, t + 0.8)
    } else {
      edge(lbx + lbw, lby + 60, tx, ty + th / 2, t)
      edge(lbx + lbw, lby + 60, tx, ty + 215 + th / 2, t + 0.8)
    }
    node(tx, ty, tw, th, 'Task 01', 'Serving requests', true)
    const x2 = P ? W - tx - tw : tx,
      y2 = P ? ty : ty + 215
    node(
      x2,
      y2,
      tw,
      th,
      scene === 1 ? 'Task stopped' : scene === 2 ? 'Replacement' : 'Task 02',
      scene === 1 ? 'Controlled exercise' : 'Serving requests',
      scene !== 1,
    )
    if (scene === 1) {
      rect(x2, y2, tw, 7, '#b77745', 2)
      text('01', P ? 150 : 1135, P ? 711 : 401, P ? 76 : 112, '#795235', 600)
      text(
        'task stays available',
        P ? 270 : 1105,
        P ? 739 : 537,
        P ? 27 : 24,
        '#795235',
      )
    } else if (scene === 2) {
      text('78s', P ? 132 : 1080, P ? 700 : 396, P ? 100 : 130, '#2c736a', 600)
      lines(
        'observed recovery',
        P ? 350 : 1090,
        P ? 737 : 548,
        P ? 410 : 345,
        P ? 30 : 25,
        '#3a655f',
        500,
      )
    } else {
      text('02', P ? 153 : 1132, P ? 708 : 397, P ? 82 : 112, '#345564', 600)
      text(
        'tasks in the design',
        P ? 285 : 1100,
        P ? 738 : 535,
        P ? 27 : 24,
        '#426778',
      )
    }
    void left
    void top
    void w
  }
  function accounts(t) {
    if (P) {
      node(210, 288, 480, 100, 'Management', 'Organization policy', true)
      const rows = [
        ['Log archive', 'Audit & security'],
        ['Workloads dev', 'Workloads prod'],
      ]
      rows.forEach((r, j) =>
        r.forEach((v, i) => {
          const x = 55 + i * 410,
            y = 466 + j * 181
          edge(W / 2, 388, x + 190, y, t + i + j)
          node(x, y, 380, 130, v, j ? 'Workload account' : 'Security account')
        }),
      )
    } else {
      node(605, 290, 390, 120, 'Management', 'Organization policy', true)
      ;[
        'Log archive',
        'Audit & security',
        'Workloads dev',
        'Workloads prod',
      ].forEach((v, i) => {
        const x = 70 + i * 375
        edge(W / 2, 410, x + 167, 550, t + i * 0.4)
        node(
          x,
          550,
          335,
          140,
          v,
          i < 2 ? 'Security account' : 'Workload account',
        )
      })
    }
  }
  return function render(film, seconds, portrait, poster = false) {
    P = portrait
    W = P ? 896 : 1600
    H = P ? 1120 : 1000
    canvas.width = W
    canvas.height = H
    accent = film.color
    c.fillStyle = '#e9eef0'
    c.fillRect(0, 0, W, H)
    const scene = Math.min(2, Math.floor(seconds / 8)),
      t = seconds - scene * 8
    const margin = P ? 52 : 80
    text(
      film.title.toUpperCase(),
      margin,
      P ? 46 : 47,
      P ? 24 : 24,
      '#3d5968',
      600,
      true,
    )
    text(
      `${String(scene + 1).padStart(2, '0')} / 03`,
      W - (P ? 174 : 200),
      P ? 46 : 47,
      22,
      '#47616f',
      400,
      true,
    )
    c.fillStyle = film.color
    c.fillRect(margin, P ? 95 : 93, P ? 65 : 84, 4)
    lines(
      headings[film.id][scene],
      margin,
      P ? 130 : 130,
      W - margin * 2,
      P ? 51 : 61,
      ink,
      600,
      1.09,
    )
    c.save()
    const enter = ease(t / 0.65)
    c.globalAlpha = 0.15 + 0.85 * enter
    c.translate(0, (1 - enter) * 20)
    if (film.id === 'cloudcostmcp' || film.id === 'llm-lint')
      terminal(film.id, scene, t)
    else if (
      film.id === 'tickethacker' ||
      (film.id === 'raizhost' && scene === 0)
    )
      screen(film.id, scene, t)
    else if (film.id === 'raizhost')
      flow(
        scene === 1
          ? [
              ['Owner editor', 'Approved content fields'],
              ['Content contract', 'Design remains in code'],
              ['Site repository', 'Reviewable change'],
            ]
          : [
              ['Content commit', 'Version history'],
              ['Site workflow', 'Build and publish'],
              ['S3 + CloudFront', 'Serve the website'],
            ],
        t,
      )
    else if (film.id === 'sre-reference-app') recovery(scene, t)
    else if (scene === 0) accounts(t)
    else
      flow(
        scene === 1
          ? [
              ['Primary region', 'us-west-2'],
              ['Pilot-light standby', 'us-east-1 · desired count 0'],
              ['Operator scales up', 'Activate standby compute'],
            ]
          : [
              ['Schedule', 'EventBridge'],
              ['Scoped role', 'Environment=dev'],
              ['Development service', 'Scale idle tasks to zero'],
            ],
        t,
      )
    c.restore()
    const bottom = P ? 942 : 856
    if (!poster)
      lines(
        captions[film.id][scene],
        margin,
        bottom,
        W - margin * 2,
        P ? 36 : 30,
        '#344e5e',
        400,
        1.35,
      )
    text(
      film.id === 'tickethacker'
        ? 'INTERFACE CAPTURE · SAMPLE DATA'
        : film.kind.toUpperCase(),
      margin,
      H - 48,
      P ? 17 : 18,
      '#4a6472',
      500,
      true,
    )
    c.fillStyle = '#b7c8d0'
    c.fillRect(0, H - 5, W, 5)
    if (!poster) {
      c.fillStyle = '#47786f'
      c.fillRect(0, H - 5, W * Math.min(1, seconds / 24), 5)
    }
    return canvas.toDataURL('image/png').split(',')[1]
  }
}
