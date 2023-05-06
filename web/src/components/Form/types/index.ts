import {
  PasswordInputProps as MantinePasswordInputProps,
  RadioGroupProps as MantineRadioGroupProps,
  RadioProps,
  SelectProps as MantineSelectProps,
  TextareaProps as MantineTextareaProps,
  TextInputProps as MantineTextInputProps,
  CheckboxGroupProps as MantineCheckboxGroupProps,
  CheckboxProps as MantineCheckboxProps,
  NumberInputProps as MantineNumberInputProps,
  MultiSelectProps as MantineMultiSelectProps,
  FileInputProps as MantineFileInputProps,
  SwitchGroupProps as MantineSwitchGroupProps,
  SwitchProps,
  GroupProps,
  StackProps,
  PinInputProps as MantinePinInputProps,
  InputWrapperBaseProps,
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
export type TextInputProps = Controlled<MantineTextInputProps>
export type PasswordInputProps = Controlled<MantinePasswordInputProps>
export type TextareaProps = Controlled<MantineTextareaProps>
export type NumberInputProps = Controlled<MantineNumberInputProps>
export type DateInputProps = Controlled<MantineDateInputProps>
export type PinInputProps = Controlled<MantinePinInputProps> &
  InputWrapperBaseProps
export type FileInputProps<T extends boolean> = Controlled<
  MantineFileInputProps<T>
>
export type SelectProps = Controlled<
  Omit<MantineSelectProps, 'data'> & {
    options: MantineSelectProps['data']
  }
>
export type MultiSelectProps = Controlled<
  Omit<MantineMultiSelectProps, 'data'> & {
    options: MantineMultiSelectProps['data']
  }
>
export type CheckboxProps = Controlled<MantineCheckboxProps>
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
