// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
interface IERC20 {
    function transfer(address recipient, uint256 amount) external returns (bool);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

contract HealthTrust {

    struct Dataset {
    string ipfsHash;
    uint8 gender;
    uint8 ageRange;
    uint8 bmiCategory;
    uint8[] chronicConditions;
    uint8[] healthMetricTypes;
    address owner;
    bool isActive;
    }


    struct Order {

    uint orderId;
    uint datasetId;
    address to; //client
    address from;// researcher
    uint256 amount;
    address tokenAddress;
    uint timestamp;
    bool completed;

    }

    mapping(uint8 => string) public genderMapping;
    mapping(uint8 => string) public ageRangeMapping;
    mapping(uint8 => string) public bmiCategoryMapping;
    mapping(uint8 => string) public chronicConditionMapping;
    mapping(uint8 => string) public healthDataMapping;


    mapping(uint=> Dataset) public datasets;

    uint public datasetCount;

    // Change the orders mapping to use a nested mapping instead of an array
    mapping(uint => mapping(uint => Order)) public orders; // datasetId => orderId => Order

    uint public orderCount;




    constructor() {
        // Initialize gender mapping
        genderMapping[0] = "Male";
        genderMapping[1] = "Female";
        genderMapping[2] = "Other";

        // Initialize age range mapping
        ageRangeMapping[0] = "18/23";
        ageRangeMapping[1] = "24/29";
        ageRangeMapping[2] = "30/35";
        ageRangeMapping[3] = "36/41";
        ageRangeMapping[4] = "42/47";
        ageRangeMapping[5] = "48/53";
        ageRangeMapping[6] = "54/59";
        ageRangeMapping[7] = "60/65";
        ageRangeMapping[8] = "66/71";
        ageRangeMapping[9] = "72/77";
        ageRangeMapping[10] = "78/83";
        ageRangeMapping[11] = "84/89";
        ageRangeMapping[12] = "90/95";
        ageRangeMapping[13] = "96/100";

        // Initialize BMI category mapping
        bmiCategoryMapping[0] = "Underweight";
        bmiCategoryMapping[1] = "Normal weight";
        bmiCategoryMapping[2] = "Overweight";
        bmiCategoryMapping[3] = "Obesity I";
        bmiCategoryMapping[4] = "Obesity II";
        bmiCategoryMapping[5] = "Obesity III";

        // Initialize chronic condition mapping
        chronicConditionMapping[0] = "Diabetes";
        chronicConditionMapping[1] = "Hypertension";
        chronicConditionMapping[2] = "Cardiovascular";
        chronicConditionMapping[3] = "Asthma";
        chronicConditionMapping[4] = "Depression/Anxiety";
        chronicConditionMapping[5] = "Other";

        // Initialize health metric type mapping
        healthDataMapping[0] = "Timestamp";
        healthDataMapping[1] = "Heart Rate";
        healthDataMapping[2] = "Respiratory Rate";
        healthDataMapping[3] = "Blood Oxygen Level";
        healthDataMapping[4] = "Body Temperature";
    }
    
    function submitDataset(
        string calldata _ipfsHash,
        uint8 gender,
        uint8 ageRange,
        uint8 bmiCategory,
        uint8[] calldata chronicConditions,
        uint8[] calldata healthMetricTypes

    ) external returns (bool) {

        require(gender <= 2, "Invalid gender");
        require(ageRange <= 13, "Invalid age range");
        require(bmiCategory <= 5, "Invalid BMI category");


        datasets[datasetCount] = Dataset({
            ipfsHash: _ipfsHash,
            gender: gender,
            ageRange: ageRange,
            bmiCategory: bmiCategory,
            chronicConditions: chronicConditions,
            healthMetricTypes: healthMetricTypes,
            owner: msg.sender,
            isActive: true
        });

        datasetCount++;
        return true;
    }


    function orderRequest(uint datasetId, uint256 amount, address tokenAddress) external returns (uint) {
        require(datasets[datasetId].isActive, "Dataset is not active");
        require(amount > 0, "Amount must be greater than 0");
        
        IERC20 token = IERC20(tokenAddress);
        require(token.transferFrom(msg.sender, address(this), amount), "Transfer failed");

        // Store order directly in the nested mapping
        orders[datasetId][orderCount] = Order({
            orderId: orderCount,
            datasetId: datasetId,
            from: msg.sender,
            to: datasets[datasetId].owner,
            amount: amount,
            tokenAddress: tokenAddress,
            timestamp: block.timestamp,
            completed: false
        });

        orderCount++;
        return orderCount - 1;
    }


    function validateOrder(uint datasetId, uint orderId, address rAddr, uint amount) external view returns (Order memory) {
        Order storage order = orders[datasetId][orderId];
        require(order.from == rAddr, "Unauthorized");
        require(order.amount == amount, "Incorrect amount");
        require(order.timestamp + 1 days > block.timestamp, "Order expired");
        require(datasets[datasetId].owner == order.to, "Receiver not matching dataset owner");

        return order;
    }


    // Update completeOrder function accordingly
    function completeOrder(uint datasetId, uint orderId) external {
        Order storage order = orders[datasetId][orderId];
        require(!order.completed, "Order already completed");
        
        IERC20 token = IERC20(order.tokenAddress);
        require(token.transfer(order.to, order.amount), "Transfer failed");
        
        order.completed = true;
    }

    function getDataset(uint datasetId) external view returns (Dataset memory) {

        return datasets[datasetId];

    }
    function getDatasetHash(uint datasetId) external view returns (string memory) {

        return datasets[datasetId].ipfsHash;

    }

    function getAllDatasets() external view returns (Dataset[] memory) {
        Dataset[] memory allDatasets = new Dataset[](datasetCount);
        for (uint i = 0; i < datasetCount; i++) {
            allDatasets[i] = datasets[i];
        }
        return allDatasets;
    }

}