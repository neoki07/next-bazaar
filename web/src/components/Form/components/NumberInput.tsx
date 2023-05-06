import { NumberInput as MantineNumberInput } from "@mantine/core";
import { NumberInputProps } from "../types";
import { useController } from "react-hook-form";
import { ErrorMessage } from "./ErrorMessage";
import { FC } from "react";

export const NumberInput: FC<NumberInputProps> = ({ label, name, ...rest }) => {
  const {
    field,
    fieldState: { error: fieldError },
    formState: { defaultValues },
  } = useController({ name });

  const error = fieldError ? (
    <ErrorMessage>{fieldError.message?.toString()}</ErrorMessage>
  ) : undefined;

  const { onChange, ...restField } = field;

  return (
    <MantineNumberInput
      id={name}
      label={label}
      onChange={(value) => {
        if (value === "") {
          onChange(defaultValues?.[name] ?? null);
        } else {
          onChange(value);
        }
      }}
      error={error}
      {...rest}
      {...restField}
    />
  );
};
