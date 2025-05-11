import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import "@oasisprotocol/sapphire-hardhat";
import "./tasks"

const accounts = ["ee1c48c4331066fd0c3b74d699f3801d63628797aea78ff2ade7ff350ea36839"];

const config: HardhatUserConfig = {
  solidity: "0.8.28",
  networks: {
    sapphire: {
      url: "https://sapphire.oasis.io",
      chainId: 0x5afe,
      accounts,
    },
    "sapphire-testnet": {
      url: "https://testnet.sapphire.oasis.io",
      accounts,
      chainId: 0x5aff,
    },
    "sapphire-localnet": {
      // docker run -it -p8544-8548:8544-8548 ghcr.io/oasisprotocol/sapphire-localnet
      url: "http://localhost:8545",
      chainId: 0x5afd,
      accounts,
    },
  },
  etherscan: {
    enabled: false
  },
  sourcify: {
    enabled: true,
  }
};

export default config;
