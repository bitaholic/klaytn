// SPDX-License-Identifier: MIT
//
// This contract is modified version of https://solidity-by-example.org/sending-ether/
// This is for test purpose only
//
pragma solidity ^0.8.3;

contract SendKlay {
    event Transfer(address indexed from, address indexed to, uint256 value);

    // the payable keyword in constructor enables sending Klay
    // while deploying this contract
    constructor() payable {}

    // A function to receiving Klay
    function contract_payable() public payable {}

    // Function to receive Ether. msg.data must be empty
    receive() external payable {}

    // Fallback function is called when msg.data is not empty
    fallback() external payable {}

    function contract_transfer(address payable _to, uint amount) public {
        // This function is no longer recommended for sending Ether.
        _to.transfer(amount);
        emit Transfer(_to, _to, amount);
    }

    function contract_send(address payable _to, uint amount) public {
        // Send returns a boolean value indicating success or failure.
        // This function is not recommended for sending Ether.
        bool sent = _to.send(amount);
        require(sent, "Failed to send Ether");
    }

    function contract_call(address payable _to, uint amount) public {
        // Call returns a boolean value indicating success or failure.
        // This is the current recommended method to use.
        (bool sent, bytes memory data) = _to.call{value : amount}("");
        require(sent, "Failed to send Ether");
    }

    function relay(address payable _to) public payable {
        bool sent = _to.send(msg.value);
        require(sent, "Failed to send Klay");
    }

    function payback(address payable _payback) public payable {
        bool sent = _payback.send(msg.value);
        require(sent, "Failed to payback");
    }
}

contract ProxySendKlay {
    function contract_payable() public payable {}

    function contract_transfer(address payable _calle, address payable _to, uint amount) public {
        SendKlay callee = SendKlay(_calle);
        callee.contract_transfer(_to, amount);
    }
}