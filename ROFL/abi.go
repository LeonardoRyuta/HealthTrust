package main

var (
	ABI_JSON = `[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        }
      ],
      "name": "DatasetSubmitted",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "orderId",
          "type": "uint256"
        }
      ],
      "name": "OrderCompleted",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "uint256",
          "name": "orderId",
          "type": "uint256"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "researcher",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "OrderCreated",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "orderId",
          "type": "uint256"
        }
      ],
      "name": "completeOrder",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "datasetCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "datasets",
      "outputs": [
        {
          "internalType": "string",
          "name": "ipfsHash",
          "type": "string"
        },
        {
          "internalType": "uint8",
          "name": "gender",
          "type": "uint8"
        },
        {
          "internalType": "uint8",
          "name": "ageRange",
          "type": "uint8"
        },
        {
          "internalType": "uint8",
          "name": "bmiCategory",
          "type": "uint8"
        },
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "bool",
          "name": "isActive",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getAllDatasets",
      "outputs": [
        {
          "components": [
            {
              "internalType": "string",
              "name": "ipfsHash",
              "type": "string"
            },
            {
              "internalType": "uint8",
              "name": "gender",
              "type": "uint8"
            },
            {
              "internalType": "uint8",
              "name": "ageRange",
              "type": "uint8"
            },
            {
              "internalType": "uint8",
              "name": "bmiCategory",
              "type": "uint8"
            },
            {
              "internalType": "uint8[]",
              "name": "chronicConditions",
              "type": "uint8[]"
            },
            {
              "internalType": "uint8[]",
              "name": "healthMetricTypes",
              "type": "uint8[]"
            },
            {
              "internalType": "address",
              "name": "owner",
              "type": "address"
            },
            {
              "internalType": "bool",
              "name": "isActive",
              "type": "bool"
            }
          ],
          "internalType": "struct HealthTrust.Dataset[]",
          "name": "out",
          "type": "tuple[]"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getDataset",
      "outputs": [
        {
          "components": [
            {
              "internalType": "string",
              "name": "ipfsHash",
              "type": "string"
            },
            {
              "internalType": "uint8",
              "name": "gender",
              "type": "uint8"
            },
            {
              "internalType": "uint8",
              "name": "ageRange",
              "type": "uint8"
            },
            {
              "internalType": "uint8",
              "name": "bmiCategory",
              "type": "uint8"
            },
            {
              "internalType": "uint8[]",
              "name": "chronicConditions",
              "type": "uint8[]"
            },
            {
              "internalType": "uint8[]",
              "name": "healthMetricTypes",
              "type": "uint8[]"
            },
            {
              "internalType": "address",
              "name": "owner",
              "type": "address"
            },
            {
              "internalType": "bool",
              "name": "isActive",
              "type": "bool"
            }
          ],
          "internalType": "struct HealthTrust.Dataset",
          "name": "",
          "type": "tuple"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getDatasetHash",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "orderId",
          "type": "uint256"
        }
      ],
      "name": "getStake",
      "outputs": [
        {
          "components": [
            {
              "internalType": "uint256",
              "name": "orderId",
              "type": "uint256"
            },
            {
              "internalType": "uint256",
              "name": "datasetId",
              "type": "uint256"
            },
            {
              "internalType": "address",
              "name": "researcher",
              "type": "address"
            },
            {
              "internalType": "address",
              "name": "patient",
              "type": "address"
            },
            {
              "internalType": "uint256",
              "name": "amount",
              "type": "uint256"
            },
            {
              "internalType": "address",
              "name": "tokenAddress",
              "type": "address"
            },
            {
              "internalType": "uint40",
              "name": "timestamp",
              "type": "uint40"
            },
            {
              "internalType": "bool",
              "name": "completed",
              "type": "bool"
            }
          ],
          "internalType": "struct HealthTrust.Order",
          "name": "",
          "type": "tuple"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "orderCount",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "tokenAddress",
          "type": "address"
        }
      ],
      "name": "orderRequest",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "orderId",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        }
      ],
      "name": "orderRequest",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "orderId",
          "type": "uint256"
        }
      ],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "orders",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "orderId",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "researcher",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "patient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "internalType": "address",
          "name": "tokenAddress",
          "type": "address"
        },
        {
          "internalType": "uint40",
          "name": "timestamp",
          "type": "uint40"
        },
        {
          "internalType": "bool",
          "name": "completed",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_ipfsHash",
          "type": "string"
        },
        {
          "internalType": "uint8",
          "name": "gender",
          "type": "uint8"
        },
        {
          "internalType": "uint8",
          "name": "ageRange",
          "type": "uint8"
        },
        {
          "internalType": "uint8",
          "name": "bmiCategory",
          "type": "uint8"
        },
        {
          "internalType": "uint8[]",
          "name": "chronicConditions",
          "type": "uint8[]"
        },
        {
          "internalType": "uint8[]",
          "name": "healthMetricTypes",
          "type": "uint8[]"
        }
      ],
      "name": "submitDataset",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "datasetId",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`
)
