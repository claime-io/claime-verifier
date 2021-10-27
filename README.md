# claime-verifier

## Supported Property Types

- [Discord User](#Discord%20AUser)

### To Be Supported...

- [Domain](#Domain)
- [Website](#Website)
- [Twitter Account](#Twitter%20Account)
- and more...

## Verification Methods

### Discord User

<details><summary>Claime Discord App verification</summary>
<br/>

Claims of this type are only available for [the Claime discord app](discord).
  
##### Evidence format
None.

##### Claim format

- `propertyType`: `Discord User ID`
- `propertyId`: `${your_discord_user_id}`
- `evidence`: `(blank)`
- `method`: `Claime Discord App`

example:

```json
{
  "propertyType": "Discord User ID",
  "propertyId": "000000000000000000",
  "Method": "Claime Discord App"
}
```

</details>

### Domain

<details><summary>TXT Record verification</summary>

##### Evidence format

```
example.com TXT "claime-ownership-claim=${your_address}"
```

##### Claim format

- `propertyType`: `Domain`
- `propertyId`: `${your_domain_name}`
- `evidence`: `${your_domain_name}` or blank
- `method`: `TXT` or blank

example:

```json
{
  "propertyType": "Domain",
  "propertyId": "example.com"
}
```

</details>

### Website

<details><summary>Meta tag verification</summary>

##### Evidence format

```
<meta name="claime-ownership-claim" content="${your_address}" />
```

##### Claim format

- `propertyType`: `Website`
- `propertyId`: `${your_website_url}`
- `evidence`: `${your_website_url}` or blank
- `method`: `MetaTag` or blank

example:

```json
{
  "propertyType": "Website",
  "propertyId": "example.com/page"
}
```

  </details>

### Twitter Account

<details><summary>Tweet verification</summary>

##### Evidence format

```
claime-ownership-claim=${your_address}
```

##### Claim format

- `propertyType`: `Twitter Account`
- `propertyId`: `${your_twitter_id}`
- `evidence`: `${your_tweet_id}`
- `method`: `Tweet` or blank

example:

```json
{
  "propertyType": "Twitter Account",
  "propertyId": "example_id",
  "evidence": "0000000000000000000"
}
```

</details>
