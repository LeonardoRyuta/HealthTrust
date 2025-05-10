import { expect } from "chai";
import { ethers } from "hardhat";
import { Contract } from "ethers";

describe("HealthTrust", function () {
  let HealthTrust;
  let healthTrust: any;
  let TestToken;
  let testToken: any;
  let owner;
  let researcher: any;
  let dataProvider: any;

  const mockIpfsHash = "QmTest123";
  const mockDataset = {
    gender: 0, // Male
    ageRange: 1, // 24-29
    bmiCategory: 1, // Normal weight
    chronicConditions: [0, 1], // Diabetes, Hypertension
    healthMetricTypes: [1, 2], // Heart Rate, Respiratory Rate
  };

  beforeEach(async function () {
    // Deploy test ERC20 token
    TestToken = await ethers.getContractFactory("TestToken");
    testToken = await TestToken.deploy();

    // Deploy HealthTrust contract
    HealthTrust = await ethers.getContractFactory("HealthTrust");
    healthTrust = await HealthTrust.deploy();

    [owner, researcher, dataProvider] = await ethers.getSigners();

    // Mint some tokens to researcher for testing
    await testToken.mint(researcher.address, ethers.parseEther("1000"));
  });

  describe("Dataset Management", function () {
    it("Should successfully submit a dataset", async function () {
      await healthTrust.connect(dataProvider).submitDataset(
        mockIpfsHash,
        mockDataset.gender,
        mockDataset.ageRange,
        mockDataset.bmiCategory,
        mockDataset.chronicConditions,
        mockDataset.healthMetricTypes
      );

      const dataset = await healthTrust.getDataset(0);
      expect(dataset.ipfsHash).to.equal(mockIpfsHash);
      expect(dataset.owner).to.equal(dataProvider.address);
      expect(dataset.isActive).to.be.true;
    });

    it("Should fail with invalid gender", async function () {
      await expect(
        healthTrust.connect(dataProvider).submitDataset(
          mockIpfsHash,
          99, // invalid gender
          mockDataset.ageRange,
          mockDataset.bmiCategory,
          mockDataset.chronicConditions,
          mockDataset.healthMetricTypes
        )
      ).to.be.revertedWith("Invalid gender");
    });
  });

  describe("Order Management", function () {
    beforeEach(async function () {
      // Submit a dataset first
      await healthTrust.connect(dataProvider).submitDataset(
        mockIpfsHash,
        mockDataset.gender,
        mockDataset.ageRange,
        mockDataset.bmiCategory,
        mockDataset.chronicConditions,
        mockDataset.healthMetricTypes
      );

      // Approve tokens for HealthTrust contract
      await testToken.connect(researcher).approve(healthTrust.target, ethers.parseEther("100"));
    });

    it("Should create an order request", async function () {
      const orderAmount = ethers.parseEther("10");
      await healthTrust.connect(researcher).orderRequest(0, orderAmount, testToken.target);

      const order = await healthTrust.orders(0, 0);
      expect(order.from).to.equal(researcher.address);
      expect(order.to).to.equal(dataProvider.address);
      expect(order.amount).to.equal(orderAmount);
      expect(order.completed).to.be.false;
    });

    it("Should validate an order correctly", async function () {
      const orderAmount = ethers.parseEther("10");
      

      await healthTrust.connect(researcher).orderRequest(0, orderAmount, testToken.target);
      
      // Get the original order for comparison
      const originalOrder = await healthTrust.orders(0, 0);

      
      // Validate the order using researcher signer
      const validatedOrder = await healthTrust.connect(researcher).validateOrder(
          0,
          0,
          researcher.address, // Use researcher.address directly
          orderAmount
      );
      
      // Compare against the researcher's address directly
      expect(validatedOrder.from).to.equal(researcher.address);
      expect(validatedOrder.to).to.equal(dataProvider.address);
      expect(validatedOrder.amount).to.equal(orderAmount);
    });

    it("Should complete an order and transfer tokens", async function () {
      const orderAmount = ethers.parseEther("10");
      await healthTrust.connect(researcher).orderRequest(0, orderAmount, testToken.target);

      const initialBalance = await testToken.balanceOf(dataProvider.address);
      
      await healthTrust.connect(dataProvider).completeOrder(0, 0);

      const finalBalance = await testToken.balanceOf(dataProvider.address);
      expect(finalBalance - initialBalance).to.equal(orderAmount);
    });
  });

  describe("Data Access", function () {
    beforeEach(async function () {
      await healthTrust.connect(dataProvider).submitDataset(
        mockIpfsHash,
        mockDataset.gender,
        mockDataset.ageRange,
        mockDataset.bmiCategory,
        mockDataset.chronicConditions,
        mockDataset.healthMetricTypes
      );
    });

    it("Should return correct dataset hash", async function () {
      const hash = await healthTrust.getDatasetHash(0);
      expect(hash).to.equal(mockIpfsHash);
    });

    it("Should return complete dataset information", async function () {
      const dataset = await healthTrust.getDataset(0);
      expect(dataset.ipfsHash).to.equal(mockIpfsHash);
      expect(dataset.gender).to.equal(mockDataset.gender);
      expect(dataset.ageRange).to.equal(mockDataset.ageRange);
      expect(dataset.bmiCategory).to.equal(mockDataset.bmiCategory);
    });
  });
});