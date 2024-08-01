import React from 'react';

export function useFormData() {
  const [longUrl, setLongUrl] = React.useState<string>('');
  const [alias, setAlias] = React.useState<string | undefined>();

  const onUrlChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setLongUrl(e.target.value);
  };

  const onAliasChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setAlias(e.target.value);
  };

  const onReset = () => {
    setLongUrl('');
    setAlias(undefined);
  };

  return { longUrl, alias, onUrlChange, onAliasChange, onReset };
}
