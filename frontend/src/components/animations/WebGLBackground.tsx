import { useRef, useMemo, useEffect, useState } from 'react'
import { Canvas, useFrame, useThree } from '@react-three/fiber'
import * as THREE from 'three'

const vertexShader = `
  uniform float uTime;
  uniform vec2 uMouse;
  uniform float uMouseRadius;

  attribute float aScale;
  attribute float aSpeed;
  attribute vec3 aOffset;

  varying float vAlpha;

  float noise(vec3 p) {
    return fract(sin(dot(p, vec3(12.9898, 78.233, 45.164))) * 43758.5453);
  }

  void main() {
    vec3 pos = position;

    float n = noise(pos * 0.5 + uTime * 0.1);
    pos.x += sin(uTime * aSpeed + aOffset.x) * 0.5;
    pos.y += cos(uTime * aSpeed * 0.8 + aOffset.y) * 0.5;
    pos.z += sin(uTime * aSpeed * 0.6 + aOffset.z) * 0.3;

    vec4 mvPosition = modelViewMatrix * vec4(pos, 1.0);

    vec2 screenPos = (projectionMatrix * mvPosition).xy / (projectionMatrix * mvPosition).w;
    float distToMouse = distance(screenPos, uMouse);
    float repulsion = smoothstep(uMouseRadius, 0.0, distToMouse);

    vec2 dir = normalize(screenPos - uMouse + 0.001);
    mvPosition.xy += dir * repulsion * 50.0;

    gl_Position = projectionMatrix * mvPosition;

    float size = aScale * (300.0 / -mvPosition.z);
    gl_PointSize = clamp(size, 1.0, 10.0);

    vAlpha = 0.3 + 0.4 * (1.0 - length(pos) / 10.0);
    vAlpha *= (1.0 - repulsion * 0.5);
  }
`

const fragmentShader = `
  varying float vAlpha;

  void main() {
    float dist = length(gl_PointCoord - vec2(0.5));
    if (dist > 0.5) discard;

    float alpha = smoothstep(0.5, 0.1, dist) * vAlpha;

    vec3 color = mix(
      vec3(0.0, 0.52, 0.86),
      vec3(0.0, 0.83, 1.0),
      gl_PointCoord.y
    );

    gl_FragColor = vec4(color, alpha);
  }
`

function Particles({ count = 2000 }) {
  const meshRef = useRef<THREE.Points>(null)
  const { viewport } = useThree()
  const mouseRef = useRef(new THREE.Vector2(0, 0))
  const targetMouseRef = useRef(new THREE.Vector2(0, 0))

  const [positions, scales, speeds, offsets] = useMemo(() => {
    const positions = new Float32Array(count * 3)
    const scales = new Float32Array(count)
    const speeds = new Float32Array(count)
    const offsets = new Float32Array(count * 3)

    for (let i = 0; i < count; i++) {
      const i3 = i * 3

      positions[i3] = (Math.random() - 0.5) * viewport.width * 3
      positions[i3 + 1] = (Math.random() - 0.5) * viewport.height * 3
      positions[i3 + 2] = (Math.random() - 0.5) * 10 - 5

      scales[i] = 0.5 + Math.random() * 1.5
      speeds[i] = 0.2 + Math.random() * 0.8

      offsets[i3] = Math.random() * Math.PI * 2
      offsets[i3 + 1] = Math.random() * Math.PI * 2
      offsets[i3 + 2] = Math.random() * Math.PI * 2
    }

    return [positions, scales, speeds, offsets]
  }, [count, viewport.width, viewport.height])

  const uniforms = useMemo(() => ({
    uTime: { value: 0 },
    uMouse: { value: new THREE.Vector2(0, 0) },
    uMouseRadius: { value: 0.3 }
  }), [])

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      targetMouseRef.current.x = (e.clientX / window.innerWidth) * 2 - 1
      targetMouseRef.current.y = -(e.clientY / window.innerHeight) * 2 + 1
    }

    window.addEventListener('mousemove', handleMouseMove)
    return () => window.removeEventListener('mousemove', handleMouseMove)
  }, [])

  useFrame((state) => {
    if (!meshRef.current) return

    const material = meshRef.current.material as THREE.ShaderMaterial
    material.uniforms.uTime.value = state.clock.elapsedTime

    mouseRef.current.lerp(targetMouseRef.current, 0.1)
    material.uniforms.uMouse.value.copy(mouseRef.current)
  })

  return (
    <points ref={meshRef}>
      <bufferGeometry>
        <bufferAttribute
          attach="attributes-position"
          count={count}
          array={positions}
          itemSize={3}
        />
        <bufferAttribute
          attach="attributes-aScale"
          count={count}
          array={scales}
          itemSize={1}
        />
        <bufferAttribute
          attach="attributes-aSpeed"
          count={count}
          array={speeds}
          itemSize={1}
        />
        <bufferAttribute
          attach="attributes-aOffset"
          count={count}
          array={offsets}
          itemSize={3}
        />
      </bufferGeometry>
      <shaderMaterial
        vertexShader={vertexShader}
        fragmentShader={fragmentShader}
        uniforms={uniforms}
        transparent
        depthWrite={false}
        blending={THREE.AdditiveBlending}
      />
    </points>
  )
}

export default function WebGLBackground() {
  const [isReduced, setIsReduced] = useState(false)

  useEffect(() => {
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    setIsReduced(prefersReducedMotion)
  }, [])

  if (isReduced) {
    return (
      <div className="fixed inset-0 -z-10 bg-gradient-to-br from-background via-background-secondary to-background" />
    )
  }

  return (
    <div className="fixed inset-0 -z-10">
      <Canvas
        camera={{ position: [0, 0, 10], fov: 60 }}
        gl={{
          alpha: true,
          antialias: false,
          powerPreference: 'high-performance',
        }}
        dpr={[1, 1.5]}
        style={{ background: 'transparent' }}
      >
        <color attach="background" args={['#050911']} />
        <Particles count={1500} />
      </Canvas>
    </div>
  )
}
