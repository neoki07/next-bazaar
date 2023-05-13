import { MainLayout } from "@/components/Layout";
import { useGetProduct } from "@/features/products";
import { useRouter } from "next/router";
import { FC } from "react";
import {
  Button,
  Container,
  Flex,
  Image,
  rem,
  Stack,
  Text,
} from "@mantine/core";
import { Price } from "@/components/Price";
import { NumberSelect, useForm } from "@/components/Form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";

interface ProductAreaProps {
  productId: string;
}

export const ProductArea: FC<ProductAreaProps> = ({ productId }) => {
  const { data: product, isLoading } = useGetProduct(productId);
  const schema = z.object({
    amount: z.number().min(1).max(10),
  });

  const [Form, methods] = useForm<{
    amount: number;
  }>({
    resolver: zodResolver(schema),
    defaultValues: {
      amount: 1,
    },
    onSubmit: (data) => {
      console.log(data);
    },
  });

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (product === undefined) {
    throw new Error("product is undefined");
  }

  return (
    <Stack>
      <Text size={28}>{product.name}</Text>
      <Text sx={(theme) => ({ color: theme.colors.gray[7] })}>
        {product.description}
      </Text>
      <Flex gap={rem(40)}>
        <Image src={product.imageUrl} alt={product.name} />
        <Form>
          <Stack w={rem(240)}>
            <Price price={product.price} />
            <NumberSelect
              w={rem(80)}
              label="Amount"
              name="amount"
              options={[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}
            />
            <Button type="submit" color="dark" fullWidth>
              Add Cart
            </Button>
          </Stack>
        </Form>
      </Flex>
    </Stack>
  );
};

interface ProductAreaProps {
  productId: string;
}

export default function ProductPage() {
  const router = useRouter();
  const { id } = router.query;
  if (Array.isArray(id)) {
    throw new Error("id is array:" + JSON.stringify(id));
  }

  return (
    <MainLayout>
      <Container>
        {id !== undefined && <ProductArea productId={id} />}
      </Container>
    </MainLayout>
  );
}
