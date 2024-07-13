import { Block } from 'baseui/block';
import { Card, StyledBody, StyledAction } from 'baseui/card';
import { LabelLarge } from 'baseui/typography';
import { Input } from 'baseui/input';
import { Button } from 'baseui/button';
import { getShortUrl } from '../api/api';


export function UserForm() {

  const onClick = async () => {
    console.log('Button clicked');
    const res = await getShortUrl('https://www.goofsdgle.com');
    console.log(res);
  };


  return (
    <Block
      display="flex"
      flexDirection="column"
      justifyContent="center"
      alignItems="center"
    >
      <Card overrides={{ Root: { style: { minWidth: '25rem' } } }}>
        {/* <FormControl> */}
        <StyledBody>
          <LabelLarge marginBottom="0.5rem">Super long URL</LabelLarge>
          <Block marginBottom="1rem">
            <Input size='compact' />
          </Block>
          <LabelLarge marginBottom="0.5rem">Alias (optional)</LabelLarge>
          <Block marginBottom="2rem">
            <Input size='compact' />
          </Block>
        </StyledBody>
        <StyledAction>
          <Block
            display="flex"
            flexDirection="row"
            justifyContent="flex-end"
          >
            <Button size='compact' onClick={onClick}>Make it short</Button>
          </Block>
        </StyledAction>
        {/* </FormControl> */}
      </Card>
    </Block>
  );
}