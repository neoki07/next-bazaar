import { useGetProductsId as useGetProductQuery } from "@/api/endpoints/products/products";
import { AxiosResponse } from "axios";
import { ApiProductResponse } from "@/api/model";
import { Product } from "../types";
import { transformProduct } from "@/features/products/utils/transform";

const transform = (response: AxiosResponse<ApiProductResponse>): Product => {
  return transformProduct(response.data);
};

export const useGetProduct = (id: string) => {
  return useGetProductQuery<Product>(id, {
    query: { select: transform },
  });
};
