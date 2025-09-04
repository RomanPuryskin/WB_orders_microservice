const btn = document.getElementById('findButton');
btn.addEventListener('click', findOrder);

function findOrder() {
    let order_id = document.getElementById('orderIdInput').value
    
    if (order_id === "") {
        document.getElementById('result').innerHTML = `Не введен ID`
        return
    };
   fetch(`/orders/${order_id}`)
        .then( responce => {
            // если ответ не успешный, выведем ошибку которую получили с эндпоинта
            if (responce.status !== 200) {
                responce.json().then( errorData => {
                    throw errorData;
                })
                .catch( errorData => {
                    document.getElementById('result').innerHTML = `
                    <h1> Error </h1>
                    <p>Code: ${errorData.code}</p>
                    <p>Message: ${errorData.msg}</p>`
                });
            } else {
                // выведем тело ответа
                responce.json().then( data => {
                    let itemsHTML = "";
                    if (data.items !== null ) {
                        data.items.forEach(item => {
                        itemsHTML += `
                            <div class="item">
                            <h3>Item:</h3>
                            <p>Chart ID: ${item.chrt_id}</p>
                            <p>Price: ${item.price}</p>
                            <p>RID: ${item.rid}</p>
                            <p>Sale: ${item.sale}</p>
                            <p>Size: ${item.size}</p>
                            <p>Total price: ${item.total_price}</p>
                            <p>NM ID: ${item.nm_id}</p>
                            <p>Brand: ${item.brand}</p>
                            <p>Status: ${item.status}</p>
                            </div>
                            `
                        })  
                    }

                    document.getElementById('result').innerHTML = `
                    <h1> Order info </h1>
                    <p>Order_id: ${data.order_uid}</p>
                    <p>Track_number: ${data.track_number}</p>
                    <p>Entry: ${data.entry}</p>
                    <div class="delivery">
                    <h2>Delivery</h2>
                        <p>Name: ${data.delivery.name}</p>
                        <p>Phone: ${data.delivery.phone}</p>
                        <p>Zip: ${data.delivery.zip}</p>
                        <p>City: ${data.delivery.city}</p>
                        <p>Address: ${data.delivery.address}</p>
                        <p>Region: ${data.delivery.region}</p>
                        <p>Email: ${data.delivery.email}</p>
                    </div>
                    <div class="payment">
                        <h2>Payment</h2>
                        <p>Transation: ${data.payment.transaction}</p>
                        <p>Request ID: ${data.payment.request_id}</p>
                        <p>Currency: ${data.payment.currency}</p>
                        <p>Provider: ${data.payment.provider}</p>
                        <p>Amount: ${data.payment.amount}</p>
                        <p>Payment dt: ${data.payment.payment_dt}</p>
                        <p>Bank: ${data.payment.bank}</p>
                        <p>Delivery cost: ${data.payment.delivery_cost}</p>
                        <p>Goods total: ${data.payment.goods_total}</p>
                        <p>Custom fee: ${data.payment.custom_fee}</p>
                    </div>
                    <div class="items">
                        <h2>Items</h2>
                        ${itemsHTML}
                    </div>
                    <p>Locale: ${data.locale}</p>
                    <p>Internal signature: ${data.internal_signature}</p>
                    <p>Customer ID: ${data.customer_id}</p>
                    <p>Delivery service: ${data.delivery_service}</p>
                    <p>ShardKey: ${data.shardkey}</p>
                    <p>SM ID: ${data.sm_id}</p>
                    <p>Date Created: ${data.date_created}</p>
                    <p>Oof shard: ${data.oof_shard}</p>
                    `
                });

            }   
        })
}