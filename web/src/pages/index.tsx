import { MainLayout } from "@/components/Layout";
import { useGetProducts } from "@/features/products";
import { ProductCard } from "@/features/products/components/ProductCard";
import { Container, Grid } from "@mantine/core";

export default function Home() {
  const { data, isLoading } = useGetProducts(1, 10);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <MainLayout>
      <Container size="lg">
        <Grid columns={4} gutter="xl">
          {data?.data.map((product) => (
            <Grid.Col key={product.id} span={1}>
              <ProductCard product={product} />
            </Grid.Col>
          ))}
        </Grid>
      </Container>
    </MainLayout>
  );
}
