import {
  GroupProps,
  InputWrapperBaseProps,
  CheckboxGroupProps as MantineCheckboxGroupProps,
  CheckboxProps as MantineCheckboxProps,
  FileInputProps as MantineFileInputProps,
  MultiSelectProps as MantineMultiSelectProps,
  NativeSelectProps as MantineNativeSelectProps,
  NumberInputProps as MantineNumberInputProps,
  PasswordInputProps as MantinePasswordInputProps,
  PinInputProps as MantinePinInputProps,
  RadioGroupProps as MantineRadioGroupProps,
  SelectProps as MantineSelectProps,
  SwitchGroupProps as MantineSwitchGroupProps,
  TextInputProps as MantineTextInputProps,
  TextareaProps as MantineTextareaProps,
  RadioProps,
  StackProps,
  SwitchProps,
} from '@mantine/core'
import { DateInputProps as MantineDateInputProps } from '@mantine/dates'
import { ReactNode } from 'react'

export type Option<OtherProps = {}> = {
  label: ReactNode
  value: any
} & OtherProps

export interface Options<OtherProps = {}> {
  options: Option<OtherProps>[]
}

export type Controlled<T> = { label: ReactNode; name: string } & T

export type Orientation =
  | { orientation?: 'horizontal'; orientationProps?: GroupProps }
  | { orientation?: 'vertical'; orientationProps?: StackProps }

export interface TextInputProps extends Controlled<MantineTextInputProps> {}

export interface PasswordInputProps
  extends Controlled<MantinePasswordInputProps> {}

export interface TextareaProps extends Controlled<MantineTextareaProps> {}

export interface NumberInputProps extends Controlled<MantineNumberInputProps> {}

export interface DateInputProps extends Controlled<MantineDateInputProps> {}

export interface PinInputProps
  extends Controlled<MantinePinInputProps>,
    Omit<InputWrapperBaseProps, 'error' | 'label'> {}

export interface FileInputProps<T extends boolean>
  extends Controlled<MantineFileInputProps<T>> {}

export interface SelectProps
  extends Controlled<
    Omit<MantineSelectProps, 'data'> & {
      options: MantineSelectProps['data']
    }
  > {}

export interface NumberSelectProps
  extends Controlled<
    Omit<MantineSelectProps, 'data'> & {
      options: readonly number[]
    }
  > {}

export interface NativeSelectProps
  extends Controlled<
    Omit<MantineNativeSelectProps, 'data'> & {
      options: MantineNativeSelectProps['data']
    }
  > {}

export interface NativeNumberSelectProps
  extends Controlled<
    Omit<MantineNativeSelectProps, 'data'> & {
      options: readonly number[]
    }
  > {}

export interface MultiSelectProps
  extends Controlled<
    Omit<MantineMultiSelectProps, 'data'> & {
      options: MantineMultiSelectProps['data']
    }
  > {}

export interface CheckboxProps extends Controlled<MantineCheckboxProps> {}

export type CheckboxGroupProps = Controlled<
  Omit<MantineCheckboxGroupProps, 'children'> &
    Options<MantineCheckboxProps> &
    Orientation
>

export type RadioGroupProps = Controlled<
  Omit<MantineRadioGroupProps, 'children'> & Options<RadioProps> & Orientation
>

export type SwitchGroupProps = Controlled<
  Omit<MantineSwitchGroupProps, 'children'> & Options<SwitchProps> & Orientation
>
