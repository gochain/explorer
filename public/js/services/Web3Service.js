const Web3Service = angular.module('Web3Service', [])
    .service('web3', ['config', function (config) {
        return config
            .then(cfg => {
                // Connect to Web3
                try {
                    if (typeof web3 !== 'undefined') {
                        // Use Mist/MetaMask's provider
                        web3 = new Web3(web3.currentProvider);
                    } else {
                        //console.log('No web3? You should consider trying MetaMask!')
                        // fallback - use your fallback strategy (local node / hosted node + in-dapp id mgmt / fail)
                        web3 = new Web3(new Web3.providers.HttpProvider(cfg.data.rpcUrl));
                    }
                } catch (err) {
                    //console.error(err);
                    console.error('no web3 detected');
                    ready.reject(err);
                }

                return web3;
            });
    }]);