const { expect } = require("chai"); // Chai is an assertion library

describe("Counter", function () {
  let Counter, counter;

  beforeEach(async function () {
    Counter = await ethers.getContractFactory("Counter");
    counter = await Counter.deploy(); // Deploy returns the deployed instance already
    // No need to call .deployed() on counter — it's not a function here
  });
  

  it("should start at 0", async function () {
    expect(await counter.count()).to.equal(0); // Initial value should be 0
  });

  it("should increment the counter", async function () {
    await counter.increment(); // Call increment
    expect(await counter.count()).to.equal(1); // Should now be 1
  });

  it("should decrement the counter", async function () {
    await counter.increment();  // First increment to avoid going below zero
    await counter.decrement();  // Now decrement back to 0
    expect(await counter.count()).to.equal(0);
  });

  it("should fail to decrement below 0", async function () {
    await expect(counter.decrement()).to.be.revertedWith("Can't go below 0");
  });
});
