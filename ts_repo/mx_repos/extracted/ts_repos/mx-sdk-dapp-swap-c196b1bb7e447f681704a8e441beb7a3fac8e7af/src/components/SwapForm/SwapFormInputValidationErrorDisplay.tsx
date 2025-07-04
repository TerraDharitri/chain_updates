import React from 'react';
import { FormikErrors, FormikTouched } from 'formik';

type SwapFormInputValidationErrorDisplayProps = {
  fieldName: string;
  errors: FormikErrors<object>;
  touched: FormikTouched<object>;
  className?: string;
};

export const SwapFormInputValidationErrorDisplay: React.FC<SwapFormInputValidationErrorDisplayProps> =
  ({
    fieldName,
    errors,
    touched,
    className = 'swap-form-input-validation-error-display'
  }) => {
    return (
      <small className={`d-flex flex-column flex-sm-row ${className}`}>
        {fieldName in errors && fieldName in touched && (
          <div
            className='invalid-feedback d-flex mt-2'
            data-testid={`invalid-${fieldName}`}
          >
            {errors[fieldName]}
          </div>
        )}
      </small>
    );
  };
