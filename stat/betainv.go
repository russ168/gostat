// Beta CDF and its inverse
// 
// References:
//
// Roger W. Abernathy and Robert P. Smith. "Applying Series Expansion
// to the Inverse Beta Distribution to Find Percentiles of the
// F-Distribution," ACM Transactions on Mathematical Software, volume
// 19, number 4, December 1993, pages 474-480.
//
// G.W. Hill and A.W. Davis. "Generalized asymptotic expansions of a
// Cornish-Fisher type," Annals of Mathematical Statistics, volume 39,
// number 8, August 1968, pages 1264-1273.


package stat

import (
	"fmt"
	"math"
)

func bisect(x, p, a, b, xtol, ptol float64) float64 {

	var x0, x1, px float64

	cdf := Beta_PDF(a, b)

	for math.Abs(x1-x0) > xtol {
		px = cdf(x)
		switch {
		case math.Abs(px-p) < ptol:
			return x
		case px < p:
			x0 = x
		case px > p:
			x1 = x
		}
		x = 0.5 * (x0 + x1)
	}
	return x
}

func betaContinuedFraction(α, β, x float64) float64 {

	var aa, del, res, qab, qap, qam, c, d, m2, m, acc float64
	var i int64
	const eps = 2.2204460492503131e-16
	const maxIter = 1000000000

	acc = 1e-16
	qab = α + β
	qap = α + 1.0
	qam = α - 1.0
	c = 1.0
	d = 1.0 - qab*x/qap

	if math.Abs(d) < eps {
		d = eps
	}
	d = 1.0 / d
	res = d

	for i = 1; i <= maxIter; i++ {
		m = (float64)(i)
		m2 = 2 * m
		aa = m * (β - m) * x / ((qam + m2) * (α + m2))
		d = 1.0 + aa*d
		if math.Abs(d) < eps {
			d = eps
		}
		c = 1.0 + aa/c
		if math.Abs(c) < eps {
			c = eps
		}
		d = 1.0 / d
		res *= d * c
		aa = -(α + m) * (qab + m) * x / ((α + m2) * (qap + m2))
		d = 1.0 + aa*d
		if math.Abs(d) < eps {
			d = eps
		}
		c = 1.0 + aa/c
		if math.Abs(c) < eps {
			c = eps
		}
		d = 1.0 / d
		del = d * c
		res *= del
		if math.Abs(del-1.0) < acc {
			return res
		}
	}

	panic(fmt.Sprintf("betaContinuedFraction(): α or β too big, or maxIter too small"))
	return -1.00
}
func Beta_CDF(α float64, β float64) func(x float64) float64 {
	return func(x float64) float64 {
		//func Beta_CDF(α , β , x float64) float64 {
		var y, res float64
		y = math.Exp(LnΓ(α+β) - LnΓ(α) - LnΓ(β) + α*math.Log(x) + β*math.Log(1.0-x))
		switch {
		case x == 0:
			res = 0.0
		case x == 1.0:
			res = 1.0
		case x < (α+1.0)/(α+β+2.0):
			res = y * betaContinuedFraction(α, β, x) / α
		default:
			res = 1.0 - y*betaContinuedFraction(β, α, 1.0-x)/β

		}
		return res
	}
}

// Beta_CDF_At() evaluates CDF of Beta distribution(α, β) at value x
func Beta_CDF_At(α, β, x float64) float64 {
	var res float64
	cdf := Beta_CDF(α, β)
	res = cdf(x)
	return res
}

func Beta_PDF_At(α, β, x float64) float64 {
	var res float64

		switch {
		case x == 0:
			res = 0.0
		case x == 1.0:
			res = 1.0
		default:
	pdf := Beta_PDF(α, β)
	res = pdf(x)

		}
	return res
}

// BetaInv_CDF_For() evaluates inverse CDF of Beta distribution(α, β) for probability p
func BetaInv_CDF_For(α float64, β float64, p float64) float64 {
	var res float64
	switch {
	case (p < 0.0 || p > 1.0):
		panic(fmt.Sprintf("p must be in range 0 < p < 1"))
		res = -1.00
	case α < 0.0:
		panic(fmt.Sprintf("α < 0"))
		res = -1.00
	case β < 0.0:
		panic(fmt.Sprintf("β < 0"))
		res = -1.00
	case p == 0.0:
		res = 0.0
	case p == 1.0:
		res = 1.0
	case p > 0.5:
		res = 1 - cdf_beta_Pinv(1-p, β, α)
	default:
		res = cdf_beta_Pinv(α, β, p)
	}
	return res

}

func cdf_beta_Pinv(α float64, β float64, p float64) float64 {
	var x, mean, lg_ab, lg_a, lg_b, lx, lambda, dP, phi, step, step0, step1 float64
	var n int64 = 0
	const tol = 1.4901161193847656e-08

	mean = α / (α + β)
	if p < 0.1 {
		/* small x */

		lg_ab = LnΓ(α + β)
		lg_a = LnΓ(α)
		lg_b = LnΓ(β)
		lx = (math.Log(α) + lg_a + lg_b - lg_ab + math.Log(p)) / α
		if lx <= 0 {
			x = math.Exp(lx)             /* first approximation */
			x *= math.Pow(1-x, -(β-1)/α) /* second approximation */
		} else {
			x = mean
		}

		if x > mean {
			x = mean
		}
	} else {
		/* Use expected value as first guess */
		x = mean
	}

	/* Do bisection to get closer */
	x = bisect(x, p, α, β, 0.01, 0.01)

	step0 = 999999

end:

	for math.Abs(step0) > 1e-11*x {
		dP = p - Beta_CDF_At(α, β, x)
		phi = Beta_PDF_At(α, β, x)

		if dP == 0.0 || n > 64 {
			break end
		}

		n++
		lambda = dP / math.Max(2*math.Abs(dP/x), phi)
		step0 = lambda
		step1 = -((α-1)/x - (β-1)/(1-x)) * lambda * lambda / 2
		step = step0

		if math.Abs(step1) < math.Abs(step0) {
			step += step1
		} else {
			// scale back step to a reasonable size when too large
			step *= 2 * math.Abs(step0/step1)
		}
		if x+step > 0 && x+step < 1 {
			x += step
		} else {
			x = math.Sqrt(x) * math.Sqrt(mean) // try a new starting point
		}

		if math.Abs(dP) > tol*p {
//			fmt.Println("failed at: α =",α , "  β =", β, "  p =", p) // just for testing purposes; delete this line and uncomment next one
//			panic(fmt.Sprintf("cdf_beta_Pinv() failed to converge"))
			 x=999.00; break end
		}
	}
	return x
}

