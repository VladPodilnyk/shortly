import { Client as Styletron } from 'styletron-engine-monolithic';
import { Provider as StyletronProvider } from 'styletron-react';
import { LightTheme, BaseProvider } from 'baseui';
import { Block } from 'baseui/block';

const engine = new Styletron();

export function Layout(props: { children: React.ReactNode }) {
  return (
    <StyletronProvider value={engine}>
      <BaseProvider theme={LightTheme}>
        <Block height="100%" width="100%" margin="2rem">
          {props.children}
        </Block>
      </BaseProvider>
    </StyletronProvider>
  );
}