import 'twin.macro'
import { HTMLAttributes, SVGAttributes } from 'react'

declare module 'react' {
  interface HTMLAttributes<T> extends AriaAttributes, DOMAttributes<T> {
    tw?: string
  }
  interface SVGAttributes<T> extends AriaAttributes, DOMAttributes<T> {
    tw?: string
  }
}
