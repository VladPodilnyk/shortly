import { Block } from "baseui/block"
import { Card, StyledBody, StyledAction } from "baseui/card"
import { LabelMedium } from "baseui/typography";
import { Input } from "baseui/input";
import { Button } from "baseui/button";

export function UserForm() {
  return (
    <Block
      display="flex"
      flexDirection="column"
      justifyContent="center"
      alignItems="center"
    >
      <Card>
        {/* <FormControl> */}
          <StyledBody>
            <LabelMedium>Shorten a long URL</LabelMedium>
            <Input />
            <LabelMedium>Provide your own alias</LabelMedium>
            <Input />
          </StyledBody>
          <StyledAction>
            <Button>
              Button Label
            </Button>
          </StyledAction>
        {/* </FormControl> */}
      </Card>
    </Block>
  );
}