# Starting the server #

Build the server by sourcint the "build" script, or download the "server" binary. The binary is statically linked, and should work on any x86 compatible Linux machine.

Run "server --help" for options. By default, the server runs on port 1718.

# Communicating with the server #

Communication with the server is done as a sequence of HTTP requests. Currently, there is no authentication implemented. The only security scheme is that all the entities created by the server are coded by hard to guess random keys (player id, game id and so on). For a game, two players receive different IDs, so one cannot send messages on behalf of the other.

The server supplies the following URLs:

`/connect/?player=[name]`

Connects the given player to the server. The reply consists of the key assigned to that player.

`/list/`

Lists the players connected to the server.

`/create/?owner_id=[owner_id]&guest=[guest]`

Creates a new game with two participants. The owner\_id parameter is the id of the creator (returned by /connect). The request returns only after the invitation has been accepted by the other side, and contains this player's game id.

`/wait/?owner_id=[owner_id]`

Waits for somebody to invite the player to a game. The request returns only after an invitation has been received and contains this player's game id.

`/send/?game_id=[game_id]&msg=[msg]`

Send a message to all the players in a given game.

`/receive/?game_id=[game_id]`

Receives a message sent by a player of the game.