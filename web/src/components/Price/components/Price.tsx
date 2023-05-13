import Decimal from "decimal.js";
import { FC } from "react";
import { formatMoneyFromDecimal } from "@/utils/money";
import { Text } from "@mantine/core";

interface PriceProps {
  price: Decimal;
}

export const Price: FC<PriceProps> = ({ price }) => {
  return (
    <Text size="xl" weight="bold">
      {formatMoneyFromDecimal(price)}
    </Text>
  );
};
