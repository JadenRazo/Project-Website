import { useRef, useState, useEffect } from 'react'
import { Canvas, useFrame } from '@react-three/fiber'
import { MeshTransmissionMaterial } from '@react-three/drei'
import * as THREE from 'three'

function Icosahedron() {
  const meshRef = useRef<THREE.Mesh>(null)
  const [hovered, setHovered] = useState(false)
  const [isReduced, setIsReduced] = useState(false)

  useEffect(() => {
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    setIsReduced(prefersReducedMotion)
  }, [])

  useFrame((state) => {
    if (!meshRef.current || isReduced) return

    const t = state.clock.elapsedTime

    meshRef.current.rotation.x = Math.sin(t * 0.3) * 0.2 + t * 0.1
    meshRef.current.rotation.y = Math.cos(t * 0.2) * 0.3 + t * 0.15
    meshRef.current.rotation.z = Math.sin(t * 0.1) * 0.1

    meshRef.current.position.y = Math.sin(t * 0.5) * 0.15
  })

  return (
    <mesh
      ref={meshRef}
      onPointerOver={() => setHovered(true)}
      onPointerOut={() => setHovered(false)}
      scale={hovered ? 1.1 : 1}
    >
      <icosahedronGeometry args={[1, 1]} />
      <MeshTransmissionMaterial
        backside
        samples={16}
        resolution={512}
        transmission={1}
        roughness={0.1}
        thickness={0.5}
        ior={1.5}
        chromaticAberration={0.06}
        anisotropy={0.1}
        distortion={0.2}
        distortionScale={0.3}
        temporalDistortion={0.2}
        clearcoat={1}
        attenuationDistance={0.5}
        attenuationColor="#ffffff"
        color="#0086dc"
      />
    </mesh>
  )
}

interface GlassIcosahedronProps {
  className?: string
  size?: number
}

export default function GlassIcosahedron({ className = '', size = 200 }: GlassIcosahedronProps) {
  const [isReduced, setIsReduced] = useState(false)

  useEffect(() => {
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    setIsReduced(prefersReducedMotion)
  }, [])

  if (isReduced) {
    return (
      <div
        className={`${className}`}
        style={{ width: size, height: size }}
      >
        <div className="w-full h-full rounded-full bg-gradient-to-br from-primary/20 to-accent/20 backdrop-blur-sm" />
      </div>
    )
  }

  return (
    <div
      className={`${className}`}
      style={{ width: size, height: size }}
    >
      <Canvas
        camera={{ position: [0, 0, 4], fov: 45 }}
        gl={{
          alpha: true,
          antialias: true,
          powerPreference: 'high-performance',
        }}
        dpr={[1, 2]}
        style={{ background: 'transparent' }}
      >
        <ambientLight intensity={0.5} />
        <spotLight
          position={[10, 10, 10]}
          angle={0.15}
          penumbra={1}
          intensity={1}
          castShadow
        />
        <pointLight position={[-10, -10, -10]} intensity={0.5} color="#00d4ff" />
        <pointLight position={[10, -10, 10]} intensity={0.3} color="#7c3aed" />
        <Icosahedron />
      </Canvas>
    </div>
  )
}
