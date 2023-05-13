import Decimal from "decimal.js";

export interface Product {
  id: string;
  name: string;
  description?: string;
  price: Decimal;
  stockQuantity: number;
  category: string;
  seller: string;
  imageUrl?: string;
}
