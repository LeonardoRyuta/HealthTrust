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
        address owner;
        bool    isActive;
    }

        struct DataEntry {
        string  cid;
        uint8   gender;
        uint8   ageRange;
        uint8   bmiCategory;
        uint8[] chronicConditions;
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

      mapping(uint256 => Dataset) public datasets;
    uint256 public datasetCount;

    mapping(uint256 => mapping(uint256 => Order)) public orders;
    uint256 public orderCount;

    mapping(uint256 => DataEntry) public entries;
    uint256 public entryCount;

    mapping(uint8 => uint256[]) public genderIndex;
    mapping(uint8 => uint256[]) public ageIndex;
    mapping(uint8 => uint256[]) public bmiIndex;
    mapping(uint8 => uint256[]) public illnessIndex;
    mapping(bytes32 => uint256[]) public combinedIndex;

    event DatasetSubmitted(uint256 datasetId, address owner);
    event OrderCreated(uint256 indexed datasetId, uint256 indexed orderId, address indexed researcher, uint256 amount);
    event OrderCompleted(uint256 datasetId, uint256 orderId);

    /*———————————————————
        Dataset submission
    ———————————————————*/
    function submitDataset(
        string calldata _ipfsHash,
        uint8  gender,
        uint8  ageRange,
        uint8  bmiCategory,
        uint8[] calldata chronicConditions
    ) external returns (uint256 datasetId) {
        require(gender      <= 2,  "invalid gender");
        require(ageRange    <= 13, "invalid age range");
        require(bmiCategory <= 5,  "invalid BMI");

        entries[entryCount] = DataEntry(_ipfsHash, gender, ageRange, bmiCategory, chronicConditions);

        genderIndex[gender].push(entryCount);
        ageIndex[ageRange].push(entryCount);
        bmiIndex[bmiCategory].push(entryCount);

        for (uint i = 0; i < chronicConditions.length; i++) {
            illnessIndex[chronicConditions[i]].push(entryCount);
        }

        bytes32 key = keccak256(abi.encode(gender, ageRange));
        combinedIndex[key].push(entryCount);

        entryCount++;

        Dataset storage d = datasets[datasetCount];
        d.ipfsHash          = _ipfsHash;
        d.gender            = gender;
        d.ageRange          = ageRange;
        d.bmiCategory       = bmiCategory;
        d.chronicConditions = chronicConditions;
        d.owner             = msg.sender;
        d.isActive          = true;

        emit DatasetSubmitted(datasetCount, msg.sender);
        datasetCount++;
        return datasetCount - 1;
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

    function getDatasetCount() external view returns (uint256) {
        return datasetCount;
    }

    function returnIndexedSets() external view returns (uint256[] memory, uint256[] memory, uint256[] memory, uint256[] memory,uint256[] memory) {
        return (genderIndex[0], ageIndex[0], bmiIndex[0],illnessIndex[0],combinedIndex[0]);
    }
    
        function getGenderIndex(uint8 gender) external view returns (uint256[] memory) {
        return genderIndex[gender];
    }

    function getAgeIndex(uint8 ageRange) external view returns (uint256[] memory) {
        return ageIndex[ageRange];
    }

    function getBmiIndex(uint8 bmiCategory) external view returns (uint256[] memory) {
        return bmiIndex[bmiCategory];
    }
    
    function getIllnessIndex(uint8 illnessCode) external view returns (uint256[] memory) {
        return illnessIndex[illnessCode];
    }

    function getDatasetsWith(
        uint8 gender,
        uint8 ageRange,
        uint8 bmiCategory,
        uint8 illness
    ) external view returns (DataEntry[] memory) {
        uint256[] memory result;

        if (gender != 255 && ageRange != 255 && bmiCategory != 255 && illness != 255) {
            bytes32 key = keccak256(abi.encode(gender, ageRange));
            uint256[] memory combinedSet = combinedIndex[key];
            uint256[] memory illnessSet = illnessIndex[illness];
            result = _intersect(combinedSet, illnessSet);
        } else if (gender != 255 && ageRange != 255 && bmiCategory != 255) {
            bytes32 key = keccak256(abi.encode(gender, ageRange));
            result = combinedIndex[key];
        } else if (gender != 255) {
            result = genderIndex[gender];
        } else if (ageRange != 255) {
            result = ageIndex[ageRange];
        } else if (bmiCategory != 255) {
            result = bmiIndex[bmiCategory];
        } else if (illness != 255) {
            result = illnessIndex[illness];
        }

        return _retrieveSets(result);
    }


    function _intersect(uint256[] memory a, uint256[] memory b) internal pure returns (uint256[] memory) {
        uint256[] memory temp = new uint256[](a.length);
        uint256 count = 0;

        for (uint256 i = 0; i < a.length; i++) {
            for (uint256 j = 0; j < b.length; j++) {
                if (a[i] == b[j]) {
                    temp[count] = a[i];
                    count++;
                    break;
                }
            }
        }

        uint256[] memory result = new uint256[](count);
        for (uint256 i = 0; i < count; i++) {
            result[i] = temp[i];
        }

    return result;
    }


    function _retrieveSets(uint256[] memory result) internal view returns (DataEntry[] memory) {
        DataEntry[] memory entriesResult = new DataEntry[](result.length);
        for (uint256 i = 0; i < result.length; i++) {
            entriesResult[i] = entries[result[i]];
        }
        return entriesResult;
    }
    



}
