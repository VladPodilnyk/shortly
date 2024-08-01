import React from 'react';
import { EncodeResponse, ErrorResponse, getShortUrl } from '../api/api';
import { AxiosError } from 'axios';

function parseError(err: AxiosError): ErrorResponse {
  return err.response?.data as ErrorResponse;
}

export function useRequest() {
  const [requestResult, setRequestResult] = React.useState<EncodeResponse | null>(null);
  const [validationErrors, setValidationErrors] = React.useState<ErrorResponse | null>(null);
  const [loading, setLoading] = React.useState<boolean>(false);
  const [isFatalError, setIsFatalError] = React.useState<boolean>(false);

  const onResetRequestState = () => {
    setRequestResult(null);
    setValidationErrors(null);
    setIsFatalError(false);
  };

  const makeRequest = async (url: string, alias?: string) => {
    onResetRequestState();
    setLoading(true);
    try {
      const res = await getShortUrl(url, alias);
      console.log(res);
      setRequestResult(res);
    } catch (e){
      if (e instanceof AxiosError) {
        setValidationErrors(parseError(e));
      } else {
        setIsFatalError(true);
      }
    } finally {
      setLoading(false);
    }
  };

  return {
    result: requestResult,
    loading,
    isUrlInvalid: Boolean(validationErrors?.error?.url),
    urlErrorMessage: validationErrors?.error?.url,
    isAliasInvalid: Boolean(validationErrors?.error?.alias),
    aliasErrorMessage: validationErrors?.error?.alias,
    isFatalError,
    makeRequest,
    onResetRequestState
  };
}