import { ApiProductResponse } from "@/api/model";
import { Product } from "@/features/products";
import Decimal from "decimal.js";

export function transformProduct(product: ApiProductResponse): Product {
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
}
