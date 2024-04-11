import terser from '@rollup/plugin-terser'
import { nodeResolve } from '@rollup/plugin-node-resolve'
import { wasm } from '@rollup/plugin-wasm'
import commonjs from '@rollup/plugin-commonjs'
import alias from '@rollup/plugin-alias'
import json from '@rollup/plugin-json'
import nodePolyfills from 'rollup-plugin-polyfill-node'

const aliases = alias({
  entries: [
    { find: 'crypto', replacement: 'crypto-browserify' }
  ],
})

export default {
  input: 'internal/html/static/srp.js',
  output: {
    format: 'es',
    file: 'internal/html/static/assets/srp.min.js',
    name: 'srp',
    sourcemap: true,
    plugins: [terser()]
  },
  plugins: [nodeResolve(), aliases, wasm(), commonjs(), json(), nodePolyfills()]
}
