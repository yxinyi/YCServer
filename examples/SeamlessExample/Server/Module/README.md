### 登录及进入地图

```sequence
Cli->UserMgr:MSG_C2S_Login
Cli->UserMgr:MSG_C2S_EnterMap
UserMgr->MapMgr:RPC_GetLeastLoadMap
MapMgr->UserMgr:RPC_GetLeastLoadMap
UserMgr->Map:RPC_UserEnterMap
Map->UserMgr:RPC_UserEnterMap
Map->UserMgr:RPC_MapLoadChange
UserMgr->Cli:MSG_S2C_Login
```

### 移动

```sequence
Cli->UserMgr:MSG_C2S_UserMove
UserMgr->Map:RPC_UserMove
Map->UserMgr:RPC_UserMove
UserMgr->Cli:MSG_S2C_UserMove
```

### 玩家移动并且通知视野内玩家

```sequence
Cli->UserMgr:MSG_C2S_UserMove
UserMgr->Map:RPC_UserMove
Map->OtherClis_1:MSG_S2C_EntityChange
Map->OtherClis_2:MSG_S2C_EntityChange
Map->OtherClis_3:MSG_S2C_EntityChange
Map->UserMgr:RPC_UserMove
UserMgr->Cli:MSG_S2C_UserMove
```









