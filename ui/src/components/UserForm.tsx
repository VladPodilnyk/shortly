import { Block } from 'baseui/block';
import { Card, StyledBody, StyledAction } from 'baseui/card';
import { LabelLarge } from 'baseui/typography';
import { Input } from 'baseui/input';
import { Button } from 'baseui/button';
import { Notification } from 'baseui/notification';
import { FormControl } from 'baseui/form-control';
import { useFormData } from '../hooks/useFormData';
import { useRequest } from '../hooks/useRequest';


export function UserForm() {
  const { longUrl, alias, onUrlChange, onAliasChange, onReset } = useFormData();
  const rq = useRequest();

  const onClick = async () => {
    await rq.makeRequest(longUrl, alias);
  };

  const onClear = () => {
    onReset();
    rq.onResetRequestState();
  };

  return (
    <Block
      display="flex"
      flexDirection="column"
      justifyContent="center"
      alignItems="center"
    >
      {rq.isFatalError && (
        <Notification
          overrides={{ Body: { style: { width: 'auto' } } }}
          kind='negative'
        >
          Something went wrong. Please try again later.
        </Notification>
      )}
      <Card overrides={{ Root: { style: { minWidth: '25rem' } } }}>
        <StyledBody>
          <LabelLarge marginBottom="0.5rem">Super long URL</LabelLarge>
          <Block marginBottom="1rem">
            <FormControl error={rq.urlErrorMessage}>
              <Input
                id='long-url'
                value={longUrl}
                error={rq.isUrlInvalid}
                onChange={onUrlChange}
                size='compact'
              />
            </FormControl>
          </Block>
          {!rq.result ? (
            <>
              <LabelLarge marginBottom="0.5rem">Alias (optional)</LabelLarge>
              <Block marginBottom="2rem">
                <FormControl error={rq.aliasErrorMessage}>
                  <Input
                    id='alias'
                    value={alias}
                    error={rq.isAliasInvalid}
                    onChange={onAliasChange}
                    size='compact'
                  />
                </FormControl>
              </Block>
            </>
          ) : (
            <>
              <LabelLarge marginBottom="0.5rem">Short URL</LabelLarge>
              <Block marginBottom="2rem">
                <Input
                  id='short-url'
                  value={rq.result.short_url}
                  readOnly
                />
              </Block>
            </>
          )}
        </StyledBody>
        <StyledAction>
          <Block
            display="flex"
            flexDirection="row"
            justifyContent="flex-end"
          >
            {!rq.result ? (
              <Button size='compact' onClick={onClick}>Make it short</Button>
            ): (
              <Button size='compact' onClick={onClear}>Clear</Button>
            )}
          </Block>
        </StyledAction>
      </Card>
    </Block>
  );
}