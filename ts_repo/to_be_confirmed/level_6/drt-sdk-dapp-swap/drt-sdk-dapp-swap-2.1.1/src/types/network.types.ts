export interface GraphqlErrorsResponseType {
  message: string;
  extensions: {
    code: string;
  };
}

export interface AuthorizationHeadersRequestParamsType {
  url: string;
  params?: Record<string, any>;
  body: {
    query?: string;
    operationName?: string;
    variables?: Record<string, any>;
  };
  method?: 'POST' | 'post';
}
