const { chromium } = await import(process.env.PLAYWRIGHT_MODULE || 'playwright')
import { mkdir, writeFile } from 'node:fs/promises'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
const out = resolve(dirname(fileURLToPath(import.meta.url)), 'captures')
const ticketUrl = process.env.TICKET_URL || 'http://127.0.0.1:4191'
if (!['localhost', '127.0.0.1'].includes(new URL(ticketUrl).hostname))
  throw new Error('TicketHacker capture must use an isolated local dashboard')
await mkdir(out, { recursive: true })
const browser = await chromium.launch({
  headless: true,
  args: ['--no-sandbox', '--disable-dev-shm-usage'],
})
try {
  const page = await browser.newPage({
    viewport: { width: 1280, height: 820 },
    deviceScaleFactor: 1,
  })
  await page.goto('https://raizhost.com', { waitUntil: 'networkidle' })
  await page.evaluate(() => window.scrollTo({ top: 0, behavior: 'instant' }))
  await page.screenshot({ path: `${out}/raizhost-wide.png` })
  await page.setViewportSize({ width: 412, height: 820 })
  await page.screenshot({ path: `${out}/raizhost-portrait.png` })
  await page.close()
  const p = await browser.newPage({ viewport: { width: 1280, height: 820 } })
  const user = {
    id: 'demo-agent',
    tenantId: 'demo-tenant',
    email: 'agent@example.com',
    name: 'Demo agent',
    role: 'AGENT',
  }
  await p.addInitScript((user) => {
    localStorage.setItem(
      'auth-storage',
      JSON.stringify({
        state: { user, accessToken: 'local-demo', refreshToken: 'local-demo' },
        version: 0,
      }),
    )
    localStorage.setItem('accessToken', 'local-demo')
    localStorage.setItem('darkMode', 'false')
  }, user)
  const now = new Date().toISOString()
  const tickets = [
    ['Website contact form', 'WEB_CHAT', 'HIGH'],
    ['Account access question', 'EMAIL', 'NORMAL'],
    ['Welcome to the community', 'DISCORD', 'NORMAL'],
    ['Update a billing address', 'TELEGRAM', 'LOW'],
  ].map(([subject, channel, priority], i) => ({
    id: `demo-${i + 1}`,
    tenantId: 'demo-tenant',
    subject,
    status: 'OPEN',
    priority,
    channel,
    contactId: `contact-${i + 1}`,
    assigneeId: user.id,
    assignee: user,
    contact: {
      id: `contact-${i + 1}`,
      name: [
        'Alex (sample)',
        'Sam (sample)',
        'Taylor (sample)',
        'Morgan (sample)',
      ][i],
      email: 'sample@example.com',
    },
    createdAt: now,
    updatedAt: now,
    _count: { messages: 2 },
    tags: [],
  }))
  const messages = [
    {
      id: 'm1',
      ticketId: 'demo-1',
      tenantId: 'demo-tenant',
      contactId: 'contact-1',
      contact: tickets[0].contact,
      direction: 'INBOUND',
      messageType: 'TEXT',
      contentText: 'Hi! Where can I update the contact form on my website?',
      createdAt: now,
    },
    {
      id: 'm2',
      ticketId: 'demo-1',
      tenantId: 'demo-tenant',
      senderId: user.id,
      sender: user,
      direction: 'OUTBOUND',
      messageType: 'TEXT',
      contentText:
        'Open your site settings and choose Contact. You can update the destination email there.',
      createdAt: now,
    },
  ]
  await p.route('**/api/**', async (route) => {
    const path = new URL(route.request().url()).pathname
    let value = []
    if (path === '/api/tickets') {
      const priority = new URL(route.request().url()).searchParams.get(
        'priority',
      )
      value = {
        data: priority
          ? tickets.filter((t) => t.priority === priority)
          : tickets,
        nextCursor: null,
      }
    } else if (path === '/api/tickets/bulk') {
      const data = route.request().postDataJSON()
      for (const t of tickets)
        if (data.ticketIds.includes(t.id)) Object.assign(t, data.updates)
      value = { updated: data.ticketIds.length }
    } else if (path === '/api/tickets/demo-1/messages')
      value = { data: messages, nextCursor: null }
    else if (path === '/api/tickets/demo-1') {
      if (route.request().method() === 'PATCH')
        Object.assign(tickets[0], route.request().postDataJSON())
      value = tickets[0]
    } else if (path === '/api/users') value = { data: [user], nextCursor: null }
    else if (path.endsWith('/health')) value = null
    else if (path.endsWith('unread-count')) value = { count: 0 }
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify(value),
    })
  })
  await p.route('**/socket.io/**', (route) => route.abort())
  await p.goto(`${ticketUrl}/tickets`, { waitUntil: 'networkidle' })
  await p.getByText('Website contact form', { exact: true }).waitFor()
  await p.screenshot({ path: `${out}/ticket-inbox.png` })
  await p.getByRole('button', { name: 'HIGH', exact: true }).click()
  await p.waitForTimeout(600)
  await p.screenshot({ path: `${out}/ticket-context.png` })
  await p.locator('tbody input[type="checkbox"]').first().check()
  await p
    .locator('select')
    .filter({ has: p.locator('option[value="RESOLVED"]') })
    .first()
    .selectOption('RESOLVED')
  await p.waitForTimeout(300)
  await p.screenshot({ path: `${out}/ticket-resolved.png` })
  await writeFile(
    `${out}/capture.json`,
    JSON.stringify(
      {
        capturedAt: now,
        raizhost: 'https://raizhost.com',
        ticketSource: 'https://github.com/JadenRazo/TicketHacker',
        ticketRevision:
          process.env.TICKET_REVISION || 'record the exact dashboard revision',
        sampleData: true,
        network:
          'TicketHacker API intercepted locally; no customer data or external messages',
      },
      null,
      2,
    ),
  )
  await p.close()
  console.log(
    'Captured public RaizHost and three real TicketHacker interface states.',
  )
} finally {
  await browser.close()
}
