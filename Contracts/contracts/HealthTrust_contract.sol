// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

library Errors {
    string constant DATASET_INACTIVE = "HealthTrust: dataset inactive";
    string constant BAD_AMOUNT       = "HealthTrust: amount = 0";
    string constant BAD_TRANSFER     = "HealthTrust: token transfer failed";
    string constant NOT_OWNER        = "HealthTrust: only dataset owner";
    string constant ORDER_DONE       = "HealthTrust: order already done";
    string constant ORDER_EXPIRED    = "HealthTrust: order expired";
    string constant UNAUTH           = "HealthTrust: unauthorised";
}

/*  ───────────────────────────────────────────────────────────────────────────
    HealthTrust
    ─────────────────────────────────────────────────────────────────────────── */
contract HealthTrust {
    /*———————————————————
        Data structures
    ———————————————————*/
    struct Dataset {
        string  ipfsHash;
        uint8   gender;
        uint8   ageRange;
        uint8   bmiCategory;
        uint8[] chronicConditions;
        uint8[] healthMetricTypes;
        address owner;
        bool    isActive;
    }

    struct Order {
        uint256 orderId;
        uint256 datasetId;
        address researcher;   // the payer
        address patient;      // dataset owner / payee
        uint256 amount;
        address tokenAddress;
        uint40  timestamp;    // creation time (fits in 5 bytes)
        bool    completed;
    }

    /*———————————————————
        Storage
    ———————————————————*/
    mapping(uint256 => Dataset) public datasets;          // datasetId → Dataset
    uint256 public datasetCount;

    // datasetId => orderId => Order
    mapping(uint256 => mapping(uint256 => Order)) public orders;
    uint256 public orderCount;

    /*———————————————————
        Events
    ———————————————————*/
    event DatasetSubmitted(uint256 indexed datasetId, address indexed owner);
    event OrderCreated(uint256 indexed datasetId, uint256 indexed orderId,
                       address indexed researcher, uint256 amount);
    event OrderCompleted(uint256 indexed datasetId, uint256 indexed orderId);

    /*———————————————————
        Dataset submission
    ———————————————————*/
    function submitDataset(
        string calldata _ipfsHash,
        uint8  gender,
        uint8  ageRange,
        uint8  bmiCategory,
        uint8[] calldata chronicConditions,
        uint8[] calldata healthMetricTypes
    ) external returns (uint256 datasetId) {
        require(gender      <= 2,  "invalid gender");
        require(ageRange    <= 13, "invalid age range");
        require(bmiCategory <= 5,  "invalid BMI");

        /* copy calldata arrays into storage */
        uint8[] memory condTmp  = chronicConditions;
        uint8[] memory metricTmp = healthMetricTypes;

        Dataset storage d = datasets[datasetCount];
        d.ipfsHash          = _ipfsHash;
        d.gender            = gender;
        d.ageRange          = ageRange;
        d.bmiCategory       = bmiCategory;
        d.chronicConditions = condTmp;
        d.healthMetricTypes = metricTmp;
        d.owner             = msg.sender;
        d.isActive          = true;

        emit DatasetSubmitted(datasetCount, msg.sender);
        return datasetCount++;
    }

    /*———————————————————
        Order flow
    ———————————————————*/
    function orderRequest(
        uint256 datasetId,
        uint256 amount,
        address tokenAddress
    ) external returns (uint256 orderId) {
        Dataset storage ds = datasets[datasetId];
        require(ds.isActive, Errors.DATASET_INACTIVE);
        require(amount > 0,  Errors.BAD_AMOUNT);

        IERC20 token = IERC20(tokenAddress);
        token.approve(address(this), amount);
        bool ok = token.transferFrom(msg.sender, address(this), amount);
        require(ok, Errors.BAD_TRANSFER);

        orders[datasetId][orderCount] = Order({
            orderId:     orderCount,
            datasetId:   datasetId,
            researcher:  msg.sender,
            patient:     ds.owner,
            amount:      amount,
            tokenAddress: tokenAddress,
            timestamp:   uint40(block.timestamp),
            completed:   false
        });

        emit OrderCreated(datasetId, orderCount, msg.sender, amount);
        return orderCount++;
    }

    function orderRequest(
        uint256 datasetId
    ) external payable returns (uint256 orderId) {
        require(msg.value == 1, Errors.BAD_AMOUNT);
        Dataset storage ds = datasets[datasetId];
        require(ds.isActive, Errors.DATASET_INACTIVE);

        orders[datasetId][orderCount] = Order({
            orderId:     orderCount,
            datasetId:   datasetId,
            researcher:  msg.sender,
            patient:     ds.owner,
            amount:      1,
            tokenAddress: address(0),
            timestamp:   uint40(block.timestamp),
            completed:   false
        });

        emit OrderCreated(datasetId, orderCount, msg.sender, 1);
        return orderCount++;
    }

    /** read‑only check used by your ROFL back‑end */
    function getStake(uint256 datasetId, uint256 orderId)
        external view returns (Order memory)
    {
        return orders[datasetId][orderId];
    }

    /** patient calls after off‑chain compute is done */
    // We need to check that the call was made by the ROFL back‑end
    function completeOrder(uint256 datasetId, uint256 orderId) external {
        Order storage o = orders[datasetId][orderId];
        require(!o.completed, Errors.ORDER_DONE);
        require(msg.sender == o.patient, Errors.NOT_OWNER);
        require(block.timestamp <= uint256(o.timestamp) + 1 days,
                Errors.ORDER_EXPIRED);

        IERC20 token = IERC20(o.tokenAddress);
        bool ok = token.transfer(o.patient, o.amount);
        require(ok, Errors.BAD_TRANSFER);

        o.completed = true;
        emit OrderCompleted(datasetId, orderId);
    }

    /*———————————————————
        Helpers (front‑end convenience)
    ———————————————————*/
    function getDataset(uint256 id) external view returns (Dataset memory) {
        return datasets[id];
    }
    function getDatasetHash(uint256 id) external view returns (string memory) {
        return datasets[id].ipfsHash;
    }
    function getAllDatasets() external view returns (Dataset[] memory out) {
        out = new Dataset[](datasetCount);
        for (uint256 i; i < datasetCount; ++i) out[i] = datasets[i];
    }
}
