package order

type Order struct {
	OrderUID          string   `json:"order_uid" faker:"len=20"`
	TrackNumber       string   `json:"track_number" faker:"len=20"`
	Entry             string   `json:"entry" faker:"oneof: WBIL"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale" faker:"oneof: en"`
	InternalSignature string   `json:"internal_signature" faker:"len=5"`
	CustomerID        string   `json:"customer_id" faker:"word"`
	DeliveryService   string   `json:"delivery_service" faker:"word"`
	Shardkey          string   `json:"shardkey" faker:"oneof: 9"`
	SmID              int64    `json:"sm_id" faker:"boundary_start=0, boundary_end=100"`
	DateCreated       string   `json:"date_created" faker:"date"`
	OofShard          string   `json:"oof_shard" faker:"oneof: 1"`
}

type Delivery struct {
	Name    string `json:"name" faker:"name"`
	Phone   string `json:"phone" faker:"e_164_phone_number"`
	Zip     string `json:"zip" faker:"oneof: 2639809"`
	City    string `json:"city" faker:"oneof: Kiryat Mozkin"`
	Address string `json:"address" faker:"oneof: Ploshad Mira 15"`
	Region  string `json:"region" faker:"oneof: Kraiot"`
	Email   string `json:"email" faker:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" faker:"len=20"`
	RequestID    string `json:"request_id" faker:"len=20"`
	Currency     string `json:"currency" faker:"currency"`
	Provider     string `json:"provider" faker:"oneof: wbpay"`
	Amount       int64  `json:"amount" faker:"boundary_start=100, boundary_end=10000"`
	PaymentDt    int64  `json:"payment_dt" faker:"unix_time"`
	Bank         string `json:"bank" faker:"word"`
	DeliveryCost int64  `json:"delivery_cost" faker:"boundary_start=100, boundary_end=10000"`
	GoodsTotal   int64  `json:"goods_total" faker:"boundary_start=1, boundary_end=100"`
	CustomFee    int64  `json:"custom_fee" faker:"boundary_start=0, boundary_end=10000"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id" faker:"boundary_start=100, boundary_end=10000"`
	TrackNumber string `json:"track_number" faker:"len=20"`
	Price       int64  `json:"price" faker:"boundary_start=100, boundary_end=10000"`
	Rid         string `json:"rid" faker:"len=20"`
	Name        string `json:"name" faker:"first_name"`
	Sale        int64  `json:"sale" faker:"boundary_start=0, boundary_end=100"`
	Size        string `json:"size" faker:"oneof: 0"`
	TotalPrice  int64  `json:"total_price" faker:"boundary_start=50, boundary_end=10000"`
	NmID        int64  `json:"nm_id" faker:"boundary_start=1000, boundary_end=1000000"`
	Brand       string `json:"brand" faker:"word"`
	Status      int64  `json:"status" faker:"boundary_start=0, boundary_end=500"`
}
