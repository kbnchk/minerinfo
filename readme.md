# minerinfo
## go library for getting basic parameters of cgminer-api compatible miners

### supported parameters:
- miner type
- miner model
- miner version (if exists)
- current hashrate (in MHS for last 5 minutes)
- pools list


This project does not aim to fully cover miner's api, library returns only most common parameters for monitoring purposes.

### tested on:
- Antminer S19, S19J Pro, S19 XP, L3+
- Whatsminer M30S,M30S+, M30S++, M31S, M31S++, M50S++
- Cheetah F9 (Shenzhen Yongyi) works as Unknown miner type


I can't guarantee that the library will work correctly with unverified miners, because each manufacturer uses its own implementation of the cgminer api, which can be radically different.

### usage example:
``` go
	var addr string = "192.168.1.100"
	var port uint16 = 4028
	var timeout time.Duration = 1 * time.Second

	miner, err := minerinfo.New(addr, port, timeout)
	if err != nil {
		log.Fatal(err)
	}

	mType:= miner.Type()
	mModel := miner.Model()
	mVersion := miner.Version()
	mHashrate := miner.Hashrate()
	mPools := miner.Pools()
	
```
