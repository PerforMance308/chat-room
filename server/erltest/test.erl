-module(test).
-export([start/1]).

start(N)->
		lists:foreach(fun(I)->
			spawn(fun()->
					test(I)
				  end)
		end,lists:seq(1,N)).

test(I)->
		{ok,S}=gen_tcp:connect("127.0.0.1",8080, [binary,{packet,raw},{active,true},{reuseaddr,true}]),
		send(S,I).

send(S,I)->
		Res = common_pb:encode_role_login_c2s({role_login_c2s, list_to_binary(integer_to_list(I))}),
		ByteRes = erlang:iolist_to_binary(Res),
		Data = <<1:16/unsigned,ByteRes/binary>>,
		gen_tcp:send(S, Data),
		receive
				{tcp, S, Bin}->
						<<Id:16/unsigned,Msg/binary>> = Bin,
						RecvMsg = protobuffs:decode(Msg, bytes);
				{tcp_closed, S}->
						io:format("server closed")
		end.
