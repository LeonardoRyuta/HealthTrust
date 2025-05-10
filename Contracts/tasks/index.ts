import { task } from "hardhat/config";

const scAddr = "0x2468f8AB370F922bdd17ff36FF87b64AdA89D060";


task("deploy", "Deploys the contract").setAction(async (taskArgs, hre) => {
    const { ethers } = hre;
    const [deployer] = await ethers.getSigners();
    
    console.log("Deploying contracts with the account:", deployer.address);
    
    const HealthTrust = await ethers.getContractFactory("HealthTrust");
    const healthTrust = await HealthTrust.deploy();
    
    console.log("Lock contract deployed to:", await healthTrust.getAddress());
    }
);