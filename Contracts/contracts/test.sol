// SPDX-License-Identifier: MIT
// This is the software license type — MIT is a permissive open-source license
pragma solidity ^0.8.0;
// This sets the Solidity compiler version. ^0.8.0 means "compatible with 0.8.x"

contract Counter {
    // Declare a public unsigned integer variable called 'count'
    // 'public' automatically creates a getter function
    uint public count;

    // Function to increase the value of 'count' by 1
    function increment() public {
        count += 1;
    }

    // Function to decrease the value of 'count' by 1
    // It includes a safety check to prevent the value from going below zero
    function decrement() public {
        require(count > 0, "Can't go below 0");
        count -= 1;
    }
}
