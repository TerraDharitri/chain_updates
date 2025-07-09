dharitri_sc::imports!();

pub type AddLiquidityResultType<M> =
    MultiValue3<DcdtTokenPayment<M>, DcdtTokenPayment<M>, DcdtTokenPayment<M>>;

pub type RemoveLiquidityResultType<M> = MultiValue2<DcdtTokenPayment<M>, DcdtTokenPayment<M>>;

pub type SwapTokensFixedInputResultType<M> = DcdtTokenPayment<M>;

pub type SwapTokensFixedOutputResultType<M> = MultiValue2<DcdtTokenPayment<M>, DcdtTokenPayment<M>>;
