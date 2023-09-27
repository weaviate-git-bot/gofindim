import * as config from '../src/config'

declare module 'config' {
  type Config = typeof config
  export default config
}
