import { Client as Styletron } from "styletron-engine-monolithic";
import { Provider as StyletronProvider } from 'styletron-react'
// import { LightTheme, BaseProvider } from 'baseui'
import { Block } from "baseui/block"

const engine = new Styletron()

export function Layout({ children }) {
  return (
    <StyletronProvider value={engine}>
      <BaseProvider theme={LightTheme}>
        <Block height="100vh" width="100vh" margin="2rem">
          {children}
        </Block>
      </BaseProvider>
    </StyletronProvider>
  )
}