WITH TB_PRODUCTS_FILTER AS (
    SELECT
        PRODUCT.LastCost,
        PRODUCT.RetailPrice,
        VENDOR.VendorName,
        VOUCHER_LINE.SKU,
        COLOR.ColorName,
        PRODUCT.SizeCode,
        PRODUCT_STYLE.StyleName,
        PRODUCT_STYLE.ExternalCode2,
        REPLACE(REPLACE(PRODUCT_STYLE.Desc1, CHAR(13), ' '), CHAR(10), ' ') AS Desc1,
        PRODUCT_UDF5.UDF5Description
    FROM
        _RetailData..VOUCHER_LINE
		JOIN _RetailData..PRODUCT ON PRODUCT.SKU = VOUCHER_LINE.SKU
		LEFT JOIN _RetailData..COLOR ON COLOR.ColorCode = PRODUCT.ColorCode
		JOIN _RetailData..PRODUCT_STYLE ON PRODUCT_STYLE.StyleCode = PRODUCT.StyleCode
        LEFT JOIN _RetailData..PRODUCT_UDF5 ON PRODUCT_STYLE.UDF5 = PRODUCT_UDF5.UDF5
        LEFT JOIN _RetailData..VENDOR ON PRODUCT_STYLE.BrandCode = VENDOR.VendorCode
    WHERE
		ReceiveDate IS NOT NULL AND
		VOUCHER_LINE.StatusCode = 'A'
        AND VOUCHER_LINE.TypeCode = 'V'
        AND VOUCHER_LINE.StoreNo IN (5, 50)
    GROUP BY
			PRODUCT.LastCost,
			PRODUCT.RetailPrice,
			VENDOR.VendorName,
			VOUCHER_LINE.SKU,
			COLOR.ColorName,
			PRODUCT.SizeCode,
			PRODUCT_STYLE.StyleName,
			PRODUCT_STYLE.ExternalCode2,
			REPLACE(REPLACE(PRODUCT_STYLE.Desc1, CHAR(13), ' '), CHAR(10), ' '),
			PRODUCT_UDF5.UDF5Description
)
,TB_SALES AS (
	SELECT
		TB_PRODUCTS_FILTER.SKU,
		AVG(RECEIPT_LINE.ExtRetailPriceWTax/NULLIF(RECEIPT_LINE.Qty, 0)) AS 'avg_sales_amount'
	FROM
		TB_PRODUCTS_FILTER
		LEFT JOIN _RetailData..RECEIPT_LINE ON TB_PRODUCTS_FILTER.SKU = RECEIPT_LINE.SKU
	WHERE
	    RECEIPT_LINE.StatusCode = 'A'
        AND
	    RECEIPT_LINE.SalesCode IN ('S', 'R')
        AND
            (
                RECEIPT_LINE.StoreNo IN (5,50)
            )
	GROUP BY
        TB_PRODUCTS_FILTER.SKU
)
,TB_IN_TEXTIL AS (
	-- Las Cargas Textil solo son 50
	SELECT
		TB_PRODUCTS_FILTER.SKU,
        ISNULL(SUM(ISNULL(Qty, 0)), 0) AS 'stock_in_textil',
        MAX(CAST(ReceiveDate AS DATE)) AS 'last_date_in_textil'
	FROM
        TB_PRODUCTS_FILTER
		LEFT JOIN _RetailData..VOUCHER_LINE ON TB_PRODUCTS_FILTER.SKU = VOUCHER_LINE.SKU
	WHERE
        VOUCHER_LINE.StatusCode = 'A'
        AND VOUCHER_LINE.TypeCode = 'V'
        AND VOUCHER_LINE.StoreNo = 50
	GROUP BY
		TB_PRODUCTS_FILTER.SKU
)
,TB_IN_TRICOTEX AS (
	-- Las cargas de Tricotex son de 17 y 5
	SELECT
		TB_PRODUCTS_FILTER.SKU,
        ISNULL(SUM(ISNULL(Qty, 0)), 0) AS 'stock_in_tricotex',
        MAX(CAST(ReceiveDate AS DATE)) AS 'last_date_in_tricotex'
	FROM
        TB_PRODUCTS_FILTER
		LEFT JOIN _RetailData..VOUCHER_LINE ON TB_PRODUCTS_FILTER.SKU = VOUCHER_LINE.SKU
	WHERE
        VOUCHER_LINE.StatusCode = 'A'
        AND VOUCHER_LINE.TypeCode = 'V'
        AND VOUCHER_LINE.StoreNo IN (5, 17)
	GROUP BY
		TB_PRODUCTS_FILTER.SKU
)
,TB_OUT_TEXTIL AS (
	SELECT
		SKU,
		ISNULL(SUM(ISNULL(stock_out_textil, 0)), 0) AS 'stock_out_textil',
		MAX(last_date_out_textil) AS 'last_date_out_textil'
	FROM
		(
			-- Retornos
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(Qty, 0)), 0) AS 'stock_out_textil',
				MAX(CAST(ReceiveDate AS DATE)) AS 'last_date_out_textil'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..VOUCHER_LINE ON TB_PRODUCTS_FILTER.SKU = VOUCHER_LINE.SKU
			WHERE
				VOUCHER_LINE.StatusCode = 'A'
				AND VOUCHER_LINE.TypeCode = 'R'
				AND VOUCHER_LINE.StoreNo = 50
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
			UNION ALL
			-- Ventas
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(Qty, 0)), 0) AS 'stock_out_textil',
				MAX(CAST(SalesDate AS DATE)) AS 'last_date_out_textil'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..RECEIPT_LINE ON RECEIPT_LINE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
				 RECEIPT_LINE.StatusCode = 'A'
				AND RECEIPT_LINE.SalesCode IN ('S', 'R')
				-- Tiendas antes del cambio de Raz�n Social de Abril
				AND
					(
						RECEIPT_LINE.Storeno IN (5,50)
					)
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
			UNION ALL
			-- Ajustes
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(-1 * ISNULL(ADJUSTMENT_LINE.QtyDiff, 0)), 0) AS 'stock_out_textil',
				MAX(CAST(AdjustmentDate AS DATE)) AS 'last_date_out_textil'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..ADJUSTMENT_LINE ON ADJUSTMENT_LINE.SKU = TB_PRODUCTS_FILTER.SKU
				JOIN _RetailData..ADJUSTMENT ON ADJUSTMENT.AdjustmentId = ADJUSTMENT_LINE.AdjustmentId
			WHERE
				ADJUSTMENT.ReasonCode = 1
				AND
				(
					-- Tiendas
					ADJUSTMENT.StoreNo IN (5,50)
				)
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
		)
		AS ALL_OUT_TEXTIL
	GROUP BY
			SKU
)
,TB_OUT_TRICOTEX AS (
	SELECT
		SKU,
		ISNULL(SUM(ISNULL(stock_out_tricotex, 0)), 0) AS 'stock_out_tricotex',
		MAX(last_date_out_tricotex) AS 'last_date_out_tricotex'
	FROM
		(
			-- Retornos
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(Qty, 0)), 0) AS 'stock_out_tricotex',
				MAX(CAST(ReceiveDate AS DATE)) AS 'last_date_out_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..VOUCHER_LINE ON TB_PRODUCTS_FILTER.SKU = VOUCHER_LINE.SKU
			WHERE
			    VOUCHER_LINE.StatusCode = 'A'
				AND VOUCHER_LINE.TypeCode = 'R'
				AND VOUCHER_LINE.StoreNo IN (5, 17)
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
			UNION ALL
			-- Ventas
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(Qty, 0)), 0) AS 'stock_out_tricotex',
				MAX(CAST(SalesDate AS DATE)) AS 'last_date_out_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..RECEIPT_LINE ON RECEIPT_LINE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
				RECEIPT_LINE.StatusCode = 'A'
				AND RECEIPT_LINE.SalesCode IN ('S', 'R')
				AND (
						RECEIPT_LINE.Storeno IN (5,50)
					)
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
			UNION ALL
			-- Ajustes
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(-1 * ISNULL(ADJUSTMENT_LINE.QtyDiff, 0)), 0) AS 'stock_out_tricotex',
				MAX(CAST(AdjustmentDate AS DATE)) AS 'last_date_out_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..ADJUSTMENT_LINE ON ADJUSTMENT_LINE.SKU = TB_PRODUCTS_FILTER.SKU
				JOIN _RetailData..ADJUSTMENT ON ADJUSTMENT.AdjustmentId = ADJUSTMENT_LINE.AdjustmentId
			WHERE
				ADJUSTMENT.ReasonCode = 1
				AND
				(
					ADJUSTMENT.Storeno IN (5,50)
				)
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
		)
		AS ALL_OUT_TRICOTEX
	GROUP BY
			SKU
)
,TB_CURRENT_STOCK_TEXTIL AS (
	SELECT
		SKU
		,ISNULL(SUM(ISNULL(stock_current_textil, 0)), 0) AS 'stock_current_textil'
	FROM
		(
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				SUM(ISNULL(PRODUCT_STORE.OnHandQty, 0)) AS 'stock_current_textil'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..PRODUCT_STORE ON PRODUCT_STORE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
				(
					-- Tiendas
					PRODUCT_STORE.StoreNo IN (5,50)
				)
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
			UNION ALL
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				SUM(ISNULL(SLIP_LINE.OutQty, 0)) AS 'stock_current_textil'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..SLIP_LINE ON SLIP_LINE.SKU = TB_PRODUCTS_FILTER.SKU
				LEFT JOIN _RetailData..SLIP ON SLIP.SlipId = SLIP_LINE.SlipId
			WHERE
				SLIP.StatusCode = 'T'
				AND
				(
					-- Tiendas
					SLIP.StoreNo IN (5,50)
				)
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
		) AS STOCK_CURRENT_TRICOTEX
	GROUP BY
		SKU
)
, TB_CURRENT_STOCK_TRICOTEX AS (
	SELECT
		SKU
		,CAST(ISNULL(SUM(ISNULL(stock_current_tricotex, 0)), 0) AS INT) AS 'stock_current_tricotex'
	FROM
		(
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				SUM(ISNULL(PRODUCT_STORE.OnHandQty, 0)) AS 'stock_current_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..PRODUCT_STORE ON PRODUCT_STORE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
                (
                    PRODUCT_STORE.StoreNo IN (5,50)
                )
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
			UNION ALL
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				SUM(ISNULL(SLIP_LINE.OutQty, 0)) AS 'stock_current_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..SLIP_LINE ON SLIP_LINE.SKU = TB_PRODUCTS_FILTER.SKU
				LEFT JOIN _RetailData..SLIP ON SLIP.SlipId = SLIP_LINE.SlipId
			WHERE
				SLIP.StatusCode = 'T'
				AND
                (
                    -- Tiendas antes del cambio de Razn Social de Abril
                    SLIP.StoreNo IN (5,50)
                )
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
		) AS STOCK_CURRENT_TEXTIL
	GROUP BY
		SKU
)
, TB_STOCK_MOVEMENT_TEXTIL AS (
	SELECT
		SKU,
		CAST(SUM(stock_movement_textil) AS INT) AS 'stock_movement_textil'
	FROM
		(
			-- RECEIPT
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				(ISNULL(SUM(ISNULL(ABS(RECEIPT_LINE.Qty), 0)), 0) * -1) AS 'stock_movement_textil'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..RECEIPT_LINE ON RECEIPT_LINE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
				(
					-- Tiendas
					RECEIPT_LINE.StoreNo IN (5,50)
				)
				AND RECEIPT_LINE.StatusCode = 'A'
				AND RECEIPT_LINE.SalesCode IN ('S', 'R')
			GROUP BY
				TB_PRODUCTS_FILTER.SKU

			UNION ALL

			-- VOUCHER ENTRY
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(ABS(VOUCHER_LINE.Qty), 0)), 0) AS 'stock_movement_textil'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..VOUCHER_LINE ON VOUCHER_LINE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
				VOUCHER_LINE.StoreNo = 50 -- TEXTIL
				AND VOUCHER_LINE.StatusCode = 'A'
				AND VOUCHER_LINE.TypeCode = 'V'
			GROUP BY
				TB_PRODUCTS_FILTER.SKU

			UNION ALL

			-- VOUCHER EXIT
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(ABS(VOUCHER_LINE.Qty), 0)), 0) AS 'stock_movement_textil'
			FROM
				TB_PRODUCTS_FILTER
				JOIN _RetailData..VOUCHER_LINE ON VOUCHER_LINE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
				VOUCHER_LINE.StoreNo = 50 -- TEXTIL
				AND VOUCHER_LINE.StatusCode = 'A'
				AND VOUCHER_LINE.TypeCode = 'R'
			GROUP BY
				TB_PRODUCTS_FILTER.SKU

			UNION ALL

			-- ADJUSTMENT
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(ADJUSTMENT_LINE.QtyDiff, 0)), 0) AS 'stock_movement_textil'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..ADJUSTMENT_LINE ON ADJUSTMENT_LINE.SKU = TB_PRODUCTS_FILTER.SKU
				LEFT JOIN _RetailData..ADJUSTMENT ON ADJUSTMENT.AdjustmentId = ADJUSTMENT_LINE.AdjustmentId
			WHERE
				(
					ADJUSTMENT.StoreNo IN (5,50)
				)
				AND ADJUSTMENT.ReasonCode = 1
				AND ADJUSTMENT.StatusCode = 'A'
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
		) RecentDocuments_1
		GROUP BY
			SKU
)
, TB_STOCK_MOVEMENT_TRICOTEX AS (
	SELECT
		SKU,
		CAST(SUM(stock_movement_tricotex) AS INT) AS 'stock_movement_tricotex'
	FROM
		(
			-- RECEIPT
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				(ISNULL(SUM(ISNULL(ABS(RECEIPT_LINE.Qty), 0)), 0) * -1) AS 'stock_movement_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..RECEIPT_LINE ON RECEIPT_LINE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
                (
                    RECEIPT_LINE.StoreNo IN (5,50)
                )
				AND RECEIPT_LINE.StatusCode = 'A'
				AND RECEIPT_LINE.SalesCode IN ('S', 'R')
			GROUP BY
				TB_PRODUCTS_FILTER.SKU

			UNION ALL

			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(ABS(VOUCHER_LINE.Qty), 0)), 0) AS 'stock_movement_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..VOUCHER_LINE ON VOUCHER_LINE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
				VOUCHER_LINE.StoreNo IN (17, 5)
				AND VOUCHER_LINE.StatusCode = 'A'
				AND VOUCHER_LINE.TypeCode = 'V'
			GROUP BY
				TB_PRODUCTS_FILTER.SKU

			UNION ALL
			-- VOUCHER EXIT
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(ABS(VOUCHER_LINE.Qty), 0)), 0) AS 'stock_movement_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				JOIN _RetailData..VOUCHER_LINE ON VOUCHER_LINE.SKU = TB_PRODUCTS_FILTER.SKU
			WHERE
				VOUCHER_LINE.StoreNo = 5 -- TRICOTEX
				AND VOUCHER_LINE.StatusCode = 'A'
				AND VOUCHER_LINE.TypeCode = 'R'
			GROUP BY
				TB_PRODUCTS_FILTER.SKU

			UNION ALL

			-- ADJUSTMENT
			SELECT
				TB_PRODUCTS_FILTER.SKU,
				ISNULL(SUM(ISNULL(ADJUSTMENT_LINE.QtyDiff, 0)), 0) AS 'stock_movement_tricotex'
			FROM
				TB_PRODUCTS_FILTER
				LEFT JOIN _RetailData..ADJUSTMENT_LINE ON ADJUSTMENT_LINE.SKU = TB_PRODUCTS_FILTER.SKU
				LEFT JOIN _RetailData..ADJUSTMENT ON ADJUSTMENT.AdjustmentId = ADJUSTMENT_LINE.AdjustmentId
			WHERE
                (
                    ADJUSTMENT.StoreNo IN (5,50)
                )
                AND ADJUSTMENT.ReasonCode = 1
                AND ADJUSTMENT.StatusCode = 'A'
			GROUP BY
				TB_PRODUCTS_FILTER.SKU
		) RecentDocuments_1
		GROUP BY
		SKU
)
, TB_RES_DETAIL AS (
    SELECT
        MAX(TB_IN_TEXTIL.last_date_in_textil) AS last_date_in_textil,
        MAX(TB_OUT_TEXTIL.last_date_out_textil) AS last_date_out_textil,
        MAX(TB_IN_TRICOTEX.last_date_in_tricotex) AS last_date_in_tricotex,
        MAX(TB_OUT_TRICOTEX.last_date_out_tricotex) AS last_date_out_tricotex,
        TB_PRODUCTS_FILTER.UDF5Description AS 'product_brand',
        TB_PRODUCTS_FILTER.VendorName AS 'origin_name',
        TB_PRODUCTS_FILTER.ExternalCode2 AS 'op',
        TB_PRODUCTS_FILTER.StyleName AS 'style_name',
        TB_PRODUCTS_FILTER.Desc1 AS 'description',
        TB_PRODUCTS_FILTER.RetailPrice AS 'product_price',
        ROUND(AVG(TB_SALES.avg_sales_amount), 2) AS 'avg_sales_amount',

        SUM(TB_CURRENT_STOCK_TEXTIL.stock_current_textil) AS stock_current_textil,
        SUM(TB_STOCK_MOVEMENT_TEXTIL.stock_movement_textil) AS stock_movement_textil,
        SUM(TB_IN_TEXTIL.stock_in_textil) AS stock_in_textil,
        SUM(TB_OUT_TEXTIL.stock_out_textil) AS stock_out_textil,

        SUM(TB_CURRENT_STOCK_TRICOTEX.stock_current_tricotex) AS stock_current_tricotex,
        SUM(TB_STOCK_MOVEMENT_TRICOTEX.stock_movement_tricotex) AS stock_movement_tricotex,
        SUM(TB_IN_TRICOTEX.stock_in_tricotex) AS stock_in_tricotex,
        SUM(TB_OUT_TRICOTEX.stock_out_tricotex) AS stock_out_tricotex,

        ISNULL(TB_PRODUCTS_FILTER.LastCost, 0) AS 'cost_unit'
    FROM
        TB_PRODUCTS_FILTER
        LEFT JOIN TB_SALES ON TB_PRODUCTS_FILTER.SKU = TB_SALES.SKU

        LEFT JOIN TB_IN_TEXTIL ON TB_PRODUCTS_FILTER.SKU = TB_IN_TEXTIL.SKU
        LEFT JOIN TB_OUT_TEXTIL ON TB_PRODUCTS_FILTER.SKU = TB_OUT_TEXTIL.SKU
        LEFT JOIN TB_IN_TRICOTEX ON TB_PRODUCTS_FILTER.SKU = TB_IN_TRICOTEX.SKU
        LEFT JOIN TB_OUT_TRICOTEX ON TB_PRODUCTS_FILTER.SKU = TB_OUT_TRICOTEX.SKU

        LEFT JOIN TB_CURRENT_STOCK_TEXTIL ON TB_PRODUCTS_FILTER.SKU = TB_CURRENT_STOCK_TEXTIL.SKU
        LEFT JOIN TB_STOCK_MOVEMENT_TEXTIL ON TB_PRODUCTS_FILTER.SKU = TB_STOCK_MOVEMENT_TEXTIL.SKU
        LEFT JOIN TB_CURRENT_STOCK_TRICOTEX ON TB_PRODUCTS_FILTER.SKU = TB_CURRENT_STOCK_TRICOTEX.SKU
        LEFT JOIN TB_STOCK_MOVEMENT_TRICOTEX ON TB_PRODUCTS_FILTER.SKU = TB_STOCK_MOVEMENT_TRICOTEX.SKU
    GROUP BY
        TB_PRODUCTS_FILTER.UDF5Description,
        TB_PRODUCTS_FILTER.VendorName,
        TB_PRODUCTS_FILTER.ExternalCode2,
        TB_PRODUCTS_FILTER.StyleName,
        TB_PRODUCTS_FILTER.Desc1,
        TB_PRODUCTS_FILTER.RetailPrice,
        TB_PRODUCTS_FILTER.LastCost
)
SELECT
	CASE
		WHEN last_date_out_textil IS NULL THEN last_date_in_textil
		ELSE
			CASE
				WHEN last_date_out_textil > last_date_in_textil
				THEN last_date_out_textil
				ELSE last_date_in_textil
			END
	END AS 'last_date_textil',
	CASE
		WHEN last_date_out_tricotex IS NULL THEN last_date_in_tricotex
		ELSE
			CASE
				WHEN last_date_out_tricotex > last_date_in_tricotex
				THEN last_date_out_tricotex
				ELSE last_date_in_tricotex
			END
	END AS 'last_date_tricotex',
    product_brand,
    origin_name,
	op,
	style_name,
	description,
    product_price,
    avg_sales_amount,

	(ISNULL(stock_current_textil, 0) - ISNULL(stock_movement_textil, 0)) AS 'stock_initial_textil',
    CAST(ISNULL(stock_in_textil, 0) AS INT) AS 'stock_in_textil',
    CAST(ISNULL(stock_out_textil, 0) AS INT) AS 'stock_out_textil',
    (CAST((ISNULL(stock_current_textil, 0) - ISNULL(stock_movement_textil, 0)) AS INT) + CAST(ISNULL(stock_in_textil, 0) AS INT) - CAST(COALESCE(ISNULL(stock_out_textil, 0), 0) AS INT)) AS 'stock_balance_textil',

	(ISNULL(stock_current_tricotex, 0) - ISNULL(stock_movement_tricotex, 0)) AS 'stock_initial_tricotex',
    CAST(ISNULL(stock_in_tricotex, 0) AS INT) AS 'stock_in_tricotex',
    CAST(ISNULL(stock_out_tricotex, 0) AS INT) AS 'stock_out_tricotex',
    (CAST((ISNULL(stock_current_tricotex, 0) - ISNULL(stock_movement_tricotex, 0)) AS INT) + CAST(ISNULL(stock_in_tricotex, 0) AS INT) - CAST(COALESCE(ISNULL(stock_out_tricotex, 0), 0) AS INT)) AS 'stock_balance_tricotex',

    cost_unit,

    cost_unit * (CAST(ISNULL(stock_in_textil, 0) AS INT) - CAST(COALESCE(ISNULL(stock_out_textil, 0), 0) AS INT)) AS 'cost_total_textil',
    cost_unit * (CAST(ISNULL(stock_in_tricotex, 0) AS INT) - CAST(COALESCE(ISNULL(stock_out_tricotex, 0), 0) AS INT)) AS 'cost_total_tricotex',

	(ISNULL(stock_current_textil, 0) - ISNULL(stock_movement_textil, 0)) + (ISNULL(stock_current_tricotex, 0) - ISNULL(stock_movement_tricotex, 0)) AS stock_initial_total,
	CAST(ISNULL(stock_in_textil, 0) AS INT) + CAST(ISNULL(stock_in_tricotex, 0) AS INT) AS stock_in_total,
	CAST(ISNULL(stock_out_textil, 0) AS INT) + CAST(ISNULL(stock_out_tricotex, 0) AS INT) AS stock_out_total,
	(CAST((ISNULL(stock_current_textil, 0) - ISNULL(stock_movement_textil, 0)) AS INT) + CAST(ISNULL(stock_in_textil, 0) AS INT) - CAST(COALESCE(ISNULL(stock_out_textil, 0), 0) AS INT)) + (CAST((ISNULL(stock_current_tricotex, 0) - ISNULL(stock_movement_tricotex, 0)) AS INT) + CAST(ISNULL(stock_in_tricotex, 0) AS INT) - CAST(COALESCE(ISNULL(stock_out_tricotex, 0), 0) AS INT)) AS stock_balance_total
FROM
    TB_RES_DETAIL

