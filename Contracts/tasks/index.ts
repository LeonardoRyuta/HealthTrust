import { task } from "hardhat/config";

task("deploy", "Deploys the contract").setAction(async (taskArgs, hre) => {
    const { ethers } = hre;
    const [deployer] = await ethers.getSigners();
    
    console.log("Deploying contracts with the account:", deployer.address);
    
    const HealthTrust = await ethers.getContractFactory("HealthTrust");
    const healthTrust = await HealthTrust.deploy("0x0039d2114a3943959d1eb74894ff4f4796103dd3e5");
    
    console.log("Lock contract deployed to:", await healthTrust.getAddress());
    }
);  