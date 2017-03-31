-ifndef(ROLE_LOGIN_C2S_PB_H).
-define(ROLE_LOGIN_C2S_PB_H, true).
-record(role_login_c2s, {
    account = erlang:error({required, account})
}).
-endif.

-ifndef(ROLE_LOGIN_S2C_PB_H).
-define(ROLE_LOGIN_S2C_PB_H, true).
-record(role_login_s2c, {
    
}).
-endif.

