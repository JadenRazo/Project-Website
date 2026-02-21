import { useRef, useEffect, useMemo, useState } from 'react'
import { Canvas, useFrame, useThree } from '@react-three/fiber'
import { useTexture } from '@react-three/drei'
import * as THREE from 'three'

const vertexShader = `
  varying vec2 vUv;

  void main() {
    vUv = uv;
    gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
  }
`

const fragmentShader = `
  uniform sampler2D uTexture;
  uniform float uTime;
  uniform float uVelocity;
  uniform float uHover;

  varying vec2 vUv;

  float random(vec2 st) {
    return fract(sin(dot(st.xy, vec2(12.9898, 78.233))) * 43758.5453123);
  }

  void main() {
    vec2 uv = vUv;

    float absVelocity = abs(uVelocity);

    float shearAmount = uVelocity * 0.015;
    uv.x += shearAmount * (uv.y - 0.5);

    float waveAmplitude = absVelocity * 0.008;
    uv.y += sin(uv.x * 12.0 + uTime * 3.0) * waveAmplitude;

    float chromaOffset = absVelocity * 0.002;

    vec4 texR = texture2D(uTexture, uv + vec2(chromaOffset * 0.5, 0.0));
    vec4 texG = texture2D(uTexture, uv);
    vec4 texB = texture2D(uTexture, uv - vec2(chromaOffset * 0.5, 0.0));

    vec4 color = vec4(texR.r, texG.g, texB.b, 1.0);

    float brightness = 1.0 + uHover * 0.15;
    color.rgb *= brightness;

    float grain = (random(uv + uTime * 0.1) - 0.5) * 0.02 * absVelocity;
    color.rgb += grain;

    gl_FragColor = color;
  }
`

interface DistortedImageProps {
  imageUrl: string
  velocity: number
  isHovered: boolean
}

function DistortedImage({ imageUrl, velocity, isHovered }: DistortedImageProps) {
  const meshRef = useRef<THREE.Mesh>(null)
  const materialRef = useRef<THREE.ShaderMaterial>(null)
  const { viewport } = useThree()

  const texture = useTexture(imageUrl)

  const uniforms = useMemo(() => ({
    uTexture: { value: texture },
    uTime: { value: 0 },
    uVelocity: { value: 0 },
    uHover: { value: 0 }
  }), [texture])

  const currentVelocity = useRef(0)
  const currentHover = useRef(0)

  useFrame((state) => {
    if (!materialRef.current) return

    materialRef.current.uniforms.uTime.value = state.clock.elapsedTime

    currentVelocity.current += (velocity - currentVelocity.current) * 0.1
    materialRef.current.uniforms.uVelocity.value = currentVelocity.current

    const targetHover = isHovered ? 1 : 0
    currentHover.current += (targetHover - currentHover.current) * 0.1
    materialRef.current.uniforms.uHover.value = currentHover.current
  })

  const scale = useMemo(() => {
    const aspectRatio = texture.image ? texture.image.width / texture.image.height : 16 / 9
    const height = viewport.height * 0.7
    const width = height * aspectRatio
    return [Math.min(width, viewport.width * 0.8), height, 1]
  }, [texture, viewport])

  return (
    <mesh ref={meshRef} scale={scale as [number, number, number]}>
      <planeGeometry args={[1, 1, 32, 32]} />
      <shaderMaterial
        ref={materialRef}
        vertexShader={vertexShader}
        fragmentShader={fragmentShader}
        uniforms={uniforms}
      />
    </mesh>
  )
}

interface DistortedProjectCardProps {
  imageUrl: string
  velocity: number
  isHovered: boolean
  className?: string
}

export default function DistortedProjectCard({
  imageUrl,
  velocity,
  isHovered,
  className = ''
}: DistortedProjectCardProps) {
  const [isReduced, setIsReduced] = useState(false)

  useEffect(() => {
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    setIsReduced(prefersReducedMotion)
  }, [])

  if (isReduced) {
    return (
      <div className={`${className} overflow-hidden rounded-2xl`}>
        <img
          src={imageUrl}
          alt=""
          className="w-full h-full object-cover"
        />
      </div>
    )
  }

  return (
    <div className={`${className}`}>
      <Canvas
        camera={{ position: [0, 0, 5], fov: 50 }}
        gl={{
          alpha: true,
          antialias: true,
          powerPreference: 'high-performance',
        }}
        dpr={[1, 2]}
        style={{ background: 'transparent' }}
      >
        <DistortedImage
          imageUrl={imageUrl}
          velocity={velocity}
          isHovered={isHovered}
        />
      </Canvas>
    </div>
  )
}
