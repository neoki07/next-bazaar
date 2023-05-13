import { useGetProducts as useGetProductsQuery } from "@/api/endpoints/products/products";
import { AxiosResponse } from "axios";
import {
  ApiListProductsResponse,
  ApiListProductsResponseMeta,
} from "@/api/model";
import { Product } from "../types";
import { transformProduct } from "@/features/products/utils/transform";

interface GetProductsResultData {
  meta: ApiListProductsResponseMeta;
  data: Product[];
}

const transform = (
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
      query: { select: transform },
    }
  );
};
