import * as THREE from 'three'

export const DistortionMaterial = {
  uniforms: {
    uTexture: { value: null as THREE.Texture | null },
    uTime: { value: 0 },
    uVelocity: { value: 0 },
    uProgress: { value: 0 },
    uResolution: { value: new THREE.Vector2(1, 1) },
  },

  vertexShader: `
    varying vec2 vUv;
    varying vec3 vPosition;

    void main() {
      vUv = uv;
      vPosition = position;
      gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
    }
  `,

  fragmentShader: `
    uniform sampler2D uTexture;
    uniform float uTime;
    uniform float uVelocity;
    uniform float uProgress;
    uniform vec2 uResolution;

    varying vec2 vUv;
    varying vec3 vPosition;

    float random(vec2 st) {
      return fract(sin(dot(st.xy, vec2(12.9898, 78.233))) * 43758.5453123);
    }

    float noise(vec2 st) {
      vec2 i = floor(st);
      vec2 f = fract(st);
      vec2 u = f * f * (3.0 - 2.0 * f);

      return mix(
        mix(random(i + vec2(0.0, 0.0)), random(i + vec2(1.0, 0.0)), u.x),
        mix(random(i + vec2(0.0, 1.0)), random(i + vec2(1.0, 1.0)), u.x),
        u.y
      );
    }

    void main() {
      vec2 uv = vUv;

      float absVelocity = abs(uVelocity);
      float velocitySign = sign(uVelocity);

      float shearAmount = uVelocity * 0.02;
      uv.x += shearAmount * (uv.y - 0.5);

      float waveFrequency = 10.0;
      float waveAmplitude = absVelocity * 0.01;
      uv.y += sin(uv.x * waveFrequency + uTime * 2.0) * waveAmplitude;

      float chromaOffset = absVelocity * 0.003;

      vec4 texR = texture2D(uTexture, uv + vec2(chromaOffset, 0.0));
      vec4 texG = texture2D(uTexture, uv);
      vec4 texB = texture2D(uTexture, uv - vec2(chromaOffset, 0.0));

      vec4 color = vec4(texR.r, texG.g, texB.b, texG.a);

      float grain = (random(uv + uTime) - 0.5) * 0.03 * absVelocity;
      color.rgb += grain;

      float vignette = 1.0 - smoothstep(0.4, 0.9, length(uv - 0.5));
      float velocityVignette = 1.0 - absVelocity * 0.1 * (1.0 - vignette);
      color.rgb *= velocityVignette;

      gl_FragColor = color;
    }
  `
}

export function createDistortionMaterial(texture?: THREE.Texture) {
  return new THREE.ShaderMaterial({
    uniforms: {
      uTexture: { value: texture || null },
      uTime: { value: 0 },
      uVelocity: { value: 0 },
      uProgress: { value: 0 },
      uResolution: { value: new THREE.Vector2(1, 1) },
    },
    vertexShader: DistortionMaterial.vertexShader,
    fragmentShader: DistortionMaterial.fragmentShader,
    transparent: true,
  })
}
