package main


func main() {
	// conn, err := grpc.Dial("kanbango.ru:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("WTF happened")
	}
	// Code removed for brevity

	// client := pb.NewInventoryClient(conn)

	// Note how we are calling the GetBookList method on the server
	// This is available to us through the auto-generated code
	// bookList, err := client.GetBookList(context.Background(), &pb.GetBookListRequest{})

	// log.Printf("book list: %v", bookList)
}