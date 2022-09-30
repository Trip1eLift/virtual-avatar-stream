from websocket_server import WebsocketServer
import json
from datetime import datetime

# Websocket server docs: https://github.com/Pithikos/python-websocket-server
HOST = "0.0.0.0"
PORT = 5001

def WS_Starts():
    server = WebsocketServer(host=HOST, port=PORT)

    def new_connection(client, server):
        print("New client has connected to the server")
        print(f"ID: {client['id']}, Address: {client['address']}")
        #print(client['handler'].rfile.readline().decode().strip())
        server.send_message(client, "Handshake")
        return

    def on_recieve(client, server, message):
        now = datetime.now()
        now_string = now.strftime("%Y/%m/%d %H:%M:%S")
        message = f"{now_string} client: {client['id']} has sent message: {message}"
        print(message)
        server.send_message(client, message)
        return

    def on_close(client, server):
        # needs to handle disconnection cleanup
        print(f"ID: {client['id']}, Address: {client['address']}", " has left.")
        return

    server.set_fn_new_client(new_connection)
    server.set_fn_message_received(on_recieve)
    server.set_fn_client_left(on_close)
    print("Listening on: ws://" + HOST + ":" + str(PORT))
    server.run_forever()

    return 'TERMINATE'

if __name__ == "__main__":
    WS_Starts()