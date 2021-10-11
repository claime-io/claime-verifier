# Claime Discord App

This bot can grant roles that you defined to people who have NFTs in your guild.

[Add to Discord](https://discord.com/oauth2/authorize?client_id=892321773981421568&permissions=268437504&scope=applications.commands%20bot)

## How to Use

### Granting roles
(guild administrators need to do [Setup](#Setup) in advance)

1. Guild members will receive an URL to verify ownership of NFTs in DM from the bot when they join.
2. Link the guild member's Discord User ID to Ethereum EOA by registering the Discord User ID with Claime or simply by signing.
    
    If you register with Claime, you can skip this step by automatically verifying and assigning roles when you join another guild. (Under development)
    
3. If the linked EOA is found to hold an NFT, the member will be granted the role.

### Setup
1. Add the Claime Discord App to the guild.
2. Create a role that you want to grant only to the owner of the NFT, and write down the ID.
3. Set the Claime role to a higher level than the role created in <2>.
4. Set the following information using the slash command.
    - ID of the role for NFT owners (created in <2>)
    - Contract address of NFT
    - Name of the network where NFT deployed (`Mainnet` , `Polygon` or `Rinkeby` only)


## Commands

### /set

Set a role to the holders of the NFT.

**Parameters**

All parmeters are required.

|Parameter|Type|Description|
|--|--|--|
|role_id|string| RoleID that is granted to the holders of the NFT |
|contract_address|string| The Contract address of the NFT |
|network|string|The network name where the NFT is deployed<br />`Mainnet` / `Polygon` / `Rinkeby`|

### /list

List the contract address and role ID pairs

**Parameters**

No parameters.

### /delete

Delete the settings associated with the contract address

**Parameters**

All parmeters are required.

|Parameter|Type|Description|
|--|--|--|
|contract_address|string| The Contract address of the NFT |
