import { useGetProducts as useGetProductsQuery } from "@/api/endpoints/products/products";
import { AxiosResponse } from "axios";
import {
  ApiListProductsResponse,
  ApiListProductsResponseMeta,
  ApiProductResponse,
} from "@/api/model";
import { Product } from "../types";
import Decimal from "decimal.js";

const transformProduct = (product: ApiProductResponse): Product => {
  if (
    product.id === undefined ||
    product.name === undefined ||
    product.price === undefined ||
    product.stock_quantity === undefined ||
    product.category === undefined ||
    product.seller === undefined
  ) {
    throw new Error("required fields are undefined:" + JSON.stringify(product));
  }

  return {
    id: product.id,
    name: product.name,
    description: product.description,
    price: new Decimal(product.price),
    stockQuantity: product.stock_quantity,
    category: product.category,
    seller: product.seller,
    imageUrl: product.image_url,
  };
};

interface GetProductsResultData {
  meta: ApiListProductsResponseMeta;
  data: Product[];
}

const transformProducts = (
  response: AxiosResponse<ApiListProductsResponse>
): GetProductsResultData => {
  const { data } = response;
  if (data.meta === undefined || data.data === undefined) {
    throw new Error("required fields are undefined:" + JSON.stringify(data));
  }

  return {
    meta: data.meta,
    data: data.data.map((item) => {
      return transformProduct(item);
    }),
  };
};

export const useGetProducts = (page: number, pageSize: number) => {
  return useGetProductsQuery<GetProductsResultData>(
    { page_id: page, page_size: pageSize },
    {
      query: { select: transformProducts },
    }
  );
};
