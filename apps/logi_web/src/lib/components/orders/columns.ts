type Status = 'pending' | 'in_progress' | 'completed' | 'cancelled';

type Order = {
  order_id: string;
  created_by: string;
  assigned_to: string;
  address: string;
  status: Status;
  created_at: string;
  updated_at: string;
};

type OrderItem = {
  item_id: string;
  order_id: string;
  product_name: string;
  quantity: number;
  created_at: string;
  updated_at: string;
};