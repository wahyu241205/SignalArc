// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract SignalArcMarket {
    enum MarketStatus {
        Draft,
        Open,
        Closed,
        Resolved,
        Cancelled
    }

    enum Outcome {
        None,
        Yes,
        No
    }

    string public question;
    uint256 public closeTimestamp;
    address public resolver;
    MarketStatus public status;
    Outcome public winningOutcome;

    event MarketCreated(string question, uint256 closeTimestamp, address resolver);
    event MarketClosed();
    event MarketCancelled();
    event MarketResolved(Outcome winningOutcome);

    error EmptyQuestion();
    error InvalidCloseTimestamp();
    error InvalidResolver();
    error MarketNotOpen();
    error MarketAlreadyFinalized();
    error UnauthorizedResolver();
    error InvalidOutcome();

    constructor(string memory question_, uint256 closeTimestamp_, address resolver_) {
        if (bytes(question_).length == 0) {
            revert EmptyQuestion();
        }

        if (closeTimestamp_ <= block.timestamp) {
            revert InvalidCloseTimestamp();
        }

        if (resolver_ == address(0)) {
            revert InvalidResolver();
        }

        question = question_;
        closeTimestamp = closeTimestamp_;
        resolver = resolver_;
        status = MarketStatus.Open;
        winningOutcome = Outcome.None;

        emit MarketCreated(question_, closeTimestamp_, resolver_);
    }

    function isOpen() external view returns (bool) {
        return status == MarketStatus.Open && block.timestamp < closeTimestamp;
    }

    function closeMarket() external {
        if (status != MarketStatus.Open || block.timestamp < closeTimestamp) {
            revert MarketNotOpen();
        }

        status = MarketStatus.Closed;

        emit MarketClosed();
    }

    function cancelMarket() external {
        if (msg.sender != resolver) {
            revert UnauthorizedResolver();
        }

        if (status == MarketStatus.Resolved || status == MarketStatus.Cancelled) {
            revert MarketAlreadyFinalized();
        }

        status = MarketStatus.Cancelled;

        emit MarketCancelled();
    }

    function resolve(Outcome winningOutcome_) external {
        if (msg.sender != resolver) {
            revert UnauthorizedResolver();
        }

        if (status != MarketStatus.Closed) {
            revert MarketNotOpen();
        }

        if (winningOutcome_ != Outcome.Yes && winningOutcome_ != Outcome.No) {
            revert InvalidOutcome();
        }

        status = MarketStatus.Resolved;
        winningOutcome = winningOutcome_;

        emit MarketResolved(winningOutcome_);
    }
}
